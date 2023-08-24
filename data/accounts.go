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

func generateId() int64 {
	result, _ := client.Incr(ctx, "users:count").Result()
	return result
}

func Add(newAccountData map[string]any) (int64, error) {
	if usernameExists(newAccountData["username"].(string)) {
		return 0, errors.New("name taken")
	}

	hashComputer := sha256.New()
	hashComputer.Write(([]byte)(newAccountData["password"].(string)))
	newAccountData["password"] = hashComputer.Sum(nil)

	id := generateId()
	newAccountData["id"] = id

	client.HSet(ctx, "user:"+strconv.FormatInt(id, 10), newAccountData)
	client.Set(ctx, "userid:"+newAccountData["username"].(string), id, 0)

	return id, nil
}

func usernameExists(username string) bool {
	out, _ := client.Exists(ctx, "userid:"+username).Result()
	return out == 1
}

func getId(username string) int {
	id, _ := client.Get(ctx, "userid:"+username).Result()
	out, _ := strconv.Atoi(id)
	return out
}

func GetUser(id int) map[string]string {
	out, _ := client.HGetAll(ctx, "user:"+strconv.FormatInt(int64(id), 10)).Result()
	return out
}

func LoginAttempt(username string, password string) (map[string]string, error) {
	id := getId(username)
	account := GetUser(id)
	if len(account) == 0 {
		return nil, errors.New("user not found")
	}

	hashComputer := sha256.New()
	hashComputer.Write(([]byte)(password))
	if string(hashComputer.Sum(nil)) != account["password"] {
		return nil, errors.New("incorrect password")
	}

	return account, nil
}

func UpdateUser(id int, data map[string]string) map[string]string {
	key := "user:" + strconv.FormatInt(int64(id), 10)
	client.HSet(ctx, key, data).Result()
	newData, _ := client.HGetAll(ctx, key).Result()
	return newData
}
