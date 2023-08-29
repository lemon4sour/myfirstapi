package data

import (
	"context"
	"crypto/sha256"
	"errors"
	"strconv"

	redis "github.com/redis/go-redis/v9"
)

var client *redis.Client
var ctx = context.Background()

func init() {
	client = redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
}

func generateID() (int64, error) {
	return client.Incr(ctx, "users:count").Result()
}

func encryptPassword(password string) ([]byte, error) {
	hashComputer := sha256.New()
	_, err := hashComputer.Write([]byte(password))
	if err != nil {
		return nil, err
	}
	return hashComputer.Sum(nil), nil
}

func AddUser(newAccountData map[string]any) (int64, error) {
	err := checkDuplicate(newAccountData["username"].(string))
	if err != nil {
		return 0, err
	}

	encrpytedPassword, err := encryptPassword(newAccountData["password"].(string))
	if err != nil {
		return 0, err
	}
	newAccountData["password"] = encrpytedPassword

	id, err := generateID()
	if err != nil {
		return 0, err
	}
	newAccountData["id"] = id

	_, err = client.HSet(ctx, "user:"+strconv.FormatInt(id, 10), newAccountData).Result()
	if err != nil {
		return 0, err
	}
	_, err = client.Set(ctx, "userid:"+newAccountData["username"].(string), id, 0).Result()
	if err != nil {
		return 0, err
	}

	AddScore(int(id), 0)

	return id, nil
}

func checkDuplicate(username string) error {
	out, err := client.Exists(ctx, "userid:"+username).Result()
	if err != nil {
		return err
	}
	if out == 1 {
		return errors.New("username taken")
	}
	return nil
}

func fetchID(username string) (int, error) {
	id, err := client.Get(ctx, "userid:"+username).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(id)
}

func FetchUser(id int) (map[string]string, error) {
	return client.HGetAll(ctx, "user:"+strconv.FormatInt(int64(id), 10)).Result()
}

func LoginAttempt(username string, password string) (map[string]string, error) {
	id, err := fetchID(username)
	if err != nil {
		return nil, err
	}
	account, err := FetchUser(id)
	if err != nil {
		return nil, err
	}
	if len(account) == 0 {
		return nil, errors.New("user not found")
	}

	encryptedPassword, err := encryptPassword(password)
	if err != nil {
		return nil, err
	}
	if string(encryptedPassword) != account["password"] {
		return nil, errors.New("incorrect password")
	}

	return account, nil
}

func UpdateUser(id int, data map[string]string) (map[string]string, error) {
	key := "user:" + strconv.FormatInt(int64(id), 10)
	_, err := client.HSet(ctx, key, data).Result()
	if err != nil {
		return nil, err
	}
	return client.HGetAll(ctx, key).Result()
}

func AddScore(id int, score float64) {
	client.ZIncrBy(ctx, "user:scores", score, strconv.FormatInt(int64(id), 10))
}

func FetchScore(id int) float64 {
	out, _ := client.ZScore(ctx, "user:scores", strconv.FormatInt(int64(id), 10)).Result()
	return out
}

type LeaderboardPlacement struct {
	ID    int
	Score float64
	Rank  int
}

func LeaderboardPage(page int, count int) ([]LeaderboardPlacement, error) {
	if page < 0 || count <= 0 {
		return nil, errors.New("negative number")
	}

	pageindex := page * count
	result, err := client.ZRevRangeWithScores(ctx, "user:scores", int64(pageindex), int64(pageindex+count-1)).Result()
	if err != nil {
		return nil, err
	}
	output := make([]LeaderboardPlacement, 0)
	for rank, z := range result {
		newLeaderboardPlacement := LeaderboardPlacement{}
		newLeaderboardPlacement.ID, err = strconv.Atoi(z.Member.(string))
		if err != nil {
			return nil, err
		}
		newLeaderboardPlacement.Score = z.Score
		newLeaderboardPlacement.Rank = rank + 1
		output = append(output, newLeaderboardPlacement)
	}

	return output, nil
}
