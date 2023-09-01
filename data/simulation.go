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
	newUser := User{}

	id, err := generateID()
	if err != nil {
		return 0, err
	}
	newUser.Username = "player_" + strconv.FormatInt(id, 10)
	newUser.Password = constantEncryptedPassword
	newUser.Name = generateName()
	newUser.Surname = generateName()

	uploadUser(newUser)
	return id, nil
}
