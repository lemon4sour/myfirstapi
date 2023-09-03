package data

import (
	"crypto/rand"
	"errors"
	"strconv"
)

var tokenCharacters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890+/"

var tokenLen = 8

func generateToken() (string, error) {
	arr := make([]byte, tokenLen)
	if _, err := rand.Read(arr); err != nil {
		return "", err
	}
	out := make([]byte, tokenLen)
	for i := 0; i < tokenLen; i++ {
		out[i] = tokenCharacters[63&arr[i]]
	}
	return string(out), nil
}

func CreateSession(id int64) (string, error) {
	var token string
	var tokenKey string
	idStr := strconv.FormatInt(id, 10)

	res, err := redisDB.Exists(ctx, "user:"+idStr+":token").Result()
	if err != nil {
		return "", err
	}
	if res != 0 {
		return "", errors.New("already in session")
	}

	for {
		var err error
		token, err = generateToken()
		if err != nil {
			return "", err
		}
		tokenKey = "token:" + token + ":userid"

		res, err := redisDB.Exists(ctx, tokenKey).Result()
		if err != nil {
			return "", err
		}
		if res == 0 {
			break
		}
	}
	if err := redisDB.Set(ctx, tokenKey, id, 600000000000).Err(); err != nil {
		return "", err
	}
	if err := redisDB.Set(ctx, "user:"+idStr+":token", token, 600000000000).Err(); err != nil {
		return "", err
	}
	return token, nil
}

func RemoveSession(token string) error {
	idStr, err := redisDB.Get(ctx, "token:"+token+":userid").Result()
	if err != nil {
		return err
	}
	err = redisDB.Del(ctx, "user:"+idStr+":token").Err()
	if err != nil {
		return err
	}
	return redisDB.Del(ctx, "token:"+token+":userid").Err()
}

func TokenToID(token string) (int64, error) {
	idStr, err := redisDB.Get(ctx, "token:"+token+":userid").Result()
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

func SessionExists(token string) (bool, error) {
	result, err := redisDB.Exists(ctx, "token:"+token+":userid").Result()
	if result > 0 {
		return true, err
	} else {
		return false, err
	}
}
