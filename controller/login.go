package controller

import (
	"encoding/json"
	"io"
	"loginsystem/data"
	"net/http"
	"strconv"
)

func Login(w http.ResponseWriter, r *http.Request) {
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

	if !MatchesTemplate(inputJSONMap, templateLogInRequest) {
		InvalidInputResponse(w, "invalid JSON parameters")
		return
	}

	account, err := data.LoginAttempt(inputJSONMap["username"].(string), (inputJSONMap["password"].(string)))
	if err != nil {
		InvalidInputResponse(w, err.Error())
		return
	}

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
