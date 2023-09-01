package data

import (
	"strconv"
)

type User struct {
	Username string `redis:"username"`
	Password []byte `redis:"password"`
	Name     string `redis:"name"`
	Surname  string `redis:"surname"`
	ID       int64  `redis:"id"`
}

func generateID() (int64, error) {
	return client.Incr(ctx, "users:count").Result()
}

func AddUser(u User) (int64, error) {
	err := checkDuplicate(u.Username)
	if err != nil {
		return 0, err
	}

	encrpytedPassword, err := encryptPassword(string(u.Password))
	if err != nil {
		return 0, err
	}
	u.Password = encrpytedPassword

	id, err := generateID()
	if err != nil {
		return 0, err
	}
	u.ID = id

	uploadUser(u)
	return id, nil
}

func fetchID(username string) (int64, error) {
	id, err := client.Get(ctx, "user:"+username+":id").Result()
	if err != nil {
		return 0, err
	}
	out, err := strconv.Atoi(id)
	return int64(out), err
}

func FetchUser(id int64) (map[string]string, error) {
	return client.HGetAll(ctx, "user:"+strconv.FormatInt(int64(id), 10)).Result()
}

func UpdateUser(id int, data map[string]string) (map[string]string, error) {
	key := "user:" + strconv.FormatInt(int64(id), 10)
	_, err := client.HSet(ctx, key, data).Result()
	if err != nil {
		return nil, err
	}
	return client.HGetAll(ctx, key).Result()
}
