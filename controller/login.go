package controller

import (
	"encoding/json"
	"loginsystem/data"
	"net/http"
	"strconv"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")

	inputMap, exists := getMap(w, r)
	if !exists {
		return
	}

	if !MatchesTemplate(inputMap, templateLogInRequest) {
		InvalidInput(w, "invalid JSON parameters")
		return
	}

	account, err := data.LoginAttempt(inputMap["username"].(string), (inputMap["password"].(string)))
	if err != nil {
		InvalidInput(w, err.Error())
		return
	}

	id, err := strconv.Atoi(account["id"])
	if err != nil {
		ServerError(w, err.Error())
		return
	}
	token, err := data.CreateSession(int64(id))
	if err != nil {
		ServerError(w, err.Error())
		return
	}
	w.Header().Set(apiTokenKey, token)

	output := templateLoginSuccess{}
	output.Status = true
	output.Result.ID, err = strconv.Atoi(account["id"])
	output.Result.Username = account["username"]
	if err != nil {
		ServerError(w, err.Error())
		return
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		ServerError(w, err.Error())
		return
	}
	w.Write(outputJSON)
}
