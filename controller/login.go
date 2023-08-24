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

	inputJSON, _ := io.ReadAll(r.Body)

	inputJSONMap, err := JSONtoMap(inputJSON)
	if err != nil {
		InvalidResponse(w, err.Error())
		return
	}

	if !MatchesTemplate(inputJSONMap, templateLogInRequest) {
		InvalidResponse(w, "invalid JSON parameters")
		return
	}

	account, err := data.LoginAttempt(inputJSONMap["username"].(string), (inputJSONMap["password"].(string)))
	if err != nil {
		InvalidResponse(w, err.Error())
		return
	}

	output := templateLoginSuccess{}
	output.Status = true
	output.Result.Id, _ = strconv.Atoi(account["id"])
	output.Result.Username = account["username"]

	outputJSON, _ := json.Marshal(output)
	w.Write(outputJSON)
}
