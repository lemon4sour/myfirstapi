package data

import (
	"crypto/sha256"
	"errors"
)

type account struct {
	Username string
	Password []byte
	Name     string
	Surname  string
	Id       int
}

var accounts []account = make([]account, 0)

func Add(newAccountData map[string]any) (*account, error) {
	if hasDuplicate(newAccountData["username"].(string)) {
		return nil, errors.New("name taken")
	}

	hashComputer := sha256.New()
	hashComputer.Write([]byte(newAccountData["password"].(string)))
	passwordhash := hashComputer.Sum(nil)
	NewAccount := account{
		Username: newAccountData["username"].(string),
		Password: passwordhash,
		Name:     newAccountData["name"].(string),
		Surname:  newAccountData["surname"].(string),
		Id:       len(accounts) + 1,
	}
	accounts = append(accounts, NewAccount)

	return &NewAccount, nil
}

func hasDuplicate(username string) bool {
	for _, v := range accounts {
		if v.Username == username {
			return true
		}
	}
	return false
}

func GetFromName(username string) *account {
	for _, v := range accounts {
		if v.Username == username {
			return &v
		}
	}
	return nil
}

func GetFromID(id int) *account {
	if id <= len(accounts) {
		return &accounts[id-1]
	}
	return nil
}
