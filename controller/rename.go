package controller

import (
	"encoding/json"
	"io"
	"loginsystem/data"
	"net/http"
	"reflect"
	"strconv"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")

	inputJSON, _ := io.ReadAll(r.Body)

	inputJSONMap, err := JSONtoMap(inputJSON)
	if err != nil {
		InvalidResponse(w, err.Error())
		return
	}

	username, exists := inputJSONMap["username"]
	if !exists {
		InvalidResponse(w, "invalid JSON parameters")
	}
	password, exists := inputJSONMap["password"]
	if !exists {
		InvalidResponse(w, "invalid JSON parameters")
	}

	account, err := data.LoginAttempt(username.(string), password.(string))
	if err != nil {
		InvalidResponse(w, err.Error())
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

	id, _ := strconv.Atoi(account["id"])
	account = data.UpdateUser(id, dataMap)

	output := templateUserData{}
	output.Status = true
	output.Result.Id = id
	output.Result.Username = account["username"]
	output.Result.Name = account["name"]
	output.Result.Surname = account["surname"]

	outputJSON, _ := json.Marshal(output)
	w.Write(outputJSON)
}
