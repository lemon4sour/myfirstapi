package controller

import (
	"encoding/json"
	"io"
	"loginsystem/data"
	"net/http"
	"reflect"
	"strconv"
)

func Rename(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")

	inputJSON, err := io.ReadAll(r.Body)
	if err != nil {
		InvalidInputResponse(w, err.Error())
		return
	}

	inputJSONMap, err := JSONtoMap(inputJSON)
	if err != nil {
		InvalidInputResponse(w, err.Error())
		return
	}

	username, exists := inputJSONMap["username"]
	if !exists {
		InvalidInputResponse(w, "invalid JSON parameters")
	}
	password, exists := inputJSONMap["password"]
	if !exists {
		InvalidInputResponse(w, "invalid JSON parameters")
	}

	account, err := data.LoginAttempt(username.(string), password.(string))
	if err != nil {
		InvalidInputResponse(w, err.Error())
		return
	}

	dataMap := make(map[string]string)
	name, exists := inputJSONMap["name"]
	if exists && (reflect.TypeOf(name) == reflect.TypeOf("")) {
		dataMap["name"] = name.(string)
	}
	surname, exists := inputJSONMap["surname"]
	if exists && (reflect.TypeOf(surname) == reflect.TypeOf("")) {
		dataMap["surname"] = surname.(string)
	}

	id, err := strconv.Atoi(account["id"])
	if err != nil {
		ServerError(w, err.Error())
		return
	}
	account, err = data.UpdateUser(id, dataMap)
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
