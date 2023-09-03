package controller

import (
	"encoding/json"
	"loginsystem/data"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")

	inputMap, exists := getMap(w, r)
	if !exists {
		return
	}

	if !MatchesTemplate(inputMap, templateRegisterRequest) {
		InvalidInput(w, "invalid JSON parameters")
		return
	}

	user := data.User{}
	user.Username = inputMap["username"].(string)
	user.Password = ([]byte)(inputMap["password"].(string))
	user.Name = inputMap["name"].(string)
	user.Surname = inputMap["surname"].(string)

	id, err := data.AddUser(user)
	if err != nil {
		InvalidInput(w, err.Error())
		return
	}

	token, err := data.CreateSession(id)
	if err != nil {
		ServerError(w, err.Error())
		return
	}
	w.Header().Set(apiTokenKey, token)

	output := templateRegisterSuccess{}
	output.Status = true
	output.Result.ID = int(id)
	output.Result.Username = inputMap["username"].(string)

	outputJSON, err := json.Marshal(output)
	if err != nil {
		ServerError(w, err.Error())
		return
	}
	w.Write(outputJSON)
}
