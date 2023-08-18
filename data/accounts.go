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

func generateId() string {
	result, _ := client.DBSize(ctx).Result()
	return strconv.FormatInt(result+1, 10)
}

func Add(newAccountData map[string]any) (int, error) {
	if hasDuplicate(newAccountData["username"].(string)) {
		return 0, errors.New("name taken")
	}

	hashComputer := sha256.New()
	hashComputer.Write(([]byte)(newAccountData["password"].(string)))
	newAccountData["password"] = hashComputer.Sum(nil)

	id := generateId()

	client.HSet(ctx, id, newAccountData)

	out, _ := strconv.Atoi(id)

	return out, nil
}

func hasDuplicate(username string) bool {
	keys, _, _ := client.ScanType(ctx, 0, "", 0, "hash").Result()
	for _, key := range keys {
		un, _ := client.HGet(ctx, key, "username").Result()
		if username == un {
			return true
		}
	}
	return false
}

func GetFromName(username string) (map[string]string, int) {
	keys, _, _ := client.ScanType(ctx, 0, "", 0, "hash").Result()
	for _, key := range keys {
		un, _ := client.HGet(ctx, key, "username").Result()
		if username == un {
			user, _ := client.HGetAll(ctx, key).Result()
			id, _ := strconv.Atoi(key)
			return user, id
		}
	}
	return nil, 0
}

func GetFromID(id int) map[string]string {
	out, _ := client.HGetAll(ctx, strconv.FormatInt(int64(id), 10)).Result()
	return out
}
