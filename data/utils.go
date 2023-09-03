package data

import (
	"context"
	"crypto/sha256"
	"errors"
	"strconv"

	redis "github.com/redis/go-redis/v9"
)

var redisDB *redis.Client
var ctx = context.Background()

func init() {
	redisDB = redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
}

func encryptPassword(password string) ([]byte, error) {
	hashComputer := sha256.New()
	_, err := hashComputer.Write([]byte(password))
	if err != nil {
		return nil, err
	}
	return hashComputer.Sum(nil), nil
}

func checkDuplicate(username string) error {
	out, err := redisDB.Exists(ctx, "user:"+username+":id").Result()
	if err != nil {
		return err
	}
	if out == 1 {
		return errors.New("username taken")
	}
	return nil
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

func uploadUser(u User) error {
	err := redisDB.HSet(ctx, "user:"+strconv.FormatInt(u.ID, 10), u).Err()
	if err != nil {
		return err
	}
	if err := redisDB.Set(ctx, "user:"+u.Username+":id", u.ID, 0).Err(); err != nil {
		return err
	}
	AddScore(u.ID, 0)
	return nil
}
