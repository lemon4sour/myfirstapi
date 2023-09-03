package controller

import (
	"encoding/json"
	"loginsystem/data"
	"net/http"
	"reflect"
)

func Rename(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")

	token, verified := verifyToken(w, r)
	if !verified {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	inputMap, exists := getMap(w, r)
	if !exists {
		return
	}

	/*
		username, exists := inputJSONMap["username"]
		if !exists {
			InvalidInput(w, "invalid JSON parameters")
		}
		password, exists := inputJSONMap["password"]
		if !exists {
			InvalidInput(w, "invalid JSON parameters")
		}

		account, err := data.LoginAttempt(username.(string), password.(string))
		if err != nil {
			InvalidInput(w, err.Error())
			return
		}
	*/

	id, err := data.TokenToID(token)
	if err != nil {
		InvalidInput(w, err.Error())
		return
	}

	dataMap := make(map[string]string)
	name, exists := inputMap["name"]
	if exists && (reflect.TypeOf(name) == reflect.TypeOf("")) {
		dataMap["name"] = name.(string)
	}
	surname, exists := inputMap["surname"]
	if exists && (reflect.TypeOf(surname) == reflect.TypeOf("")) {
		dataMap["surname"] = surname.(string)
	}
	account, err := data.UpdateUser(id, dataMap)
	if err != nil {
		ServerError(w, err.Error())
		return
	}

	output := templateUserData{}
	output.Status = true
	output.Result.ID = id
	output.Result.Username = account["username"]
	output.Result.Name = account["name"]
	output.Result.Surname = account["surname"]

	outputJSON, err := json.Marshal(output)
	if err != nil {
		ServerError(w, err.Error())
		return
	}
	w.Write(outputJSON)
}
