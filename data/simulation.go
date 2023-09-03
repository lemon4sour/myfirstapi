package data

import (
	"math/rand"
	"strconv"
)

var constantPassword = "1234567"
var constantEncryptedPassword, _ = encryptPassword(constantPassword)

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNAOPQRSTUVWXYZ1234567890_")

func generateName() string {
	nameLength := rand.Intn(6) + 3
	out := make([]rune, nameLength)
	for i := 0; i < nameLength; i++ {
		out[i] = charset[rand.Intn(len(charset))]
	}
	return string(out)
}

func GenerateUser() (int64, error) {
	u := User{}

	id, err := generateID()
	if err != nil {
		return 0, err
	}
	u.Username = "player_" + strconv.FormatInt(id, 10)
	u.Password = constantEncryptedPassword
	u.Name = generateName()
	u.Surname = generateName()

	if err := uploadUser(u); err != nil {
		return 0, err
	}
	return id, nil
}
