package controller

import (
	"encoding/json"
	"io"
	"loginsystem/data"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
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

	if !MatchesTemplate(inputJSONMap, templateRegisterRequest) {
		InvalidInputResponse(w, "invalid JSON parameters")
		return
	}

	user := data.User{}
	user.Username = inputJSONMap["username"].(string)
	user.Password = ([]byte)(inputJSONMap["password"].(string))
	user.Name = inputJSONMap["name"].(string)
	user.Surname = inputJSONMap["surname"].(string)

	id, err := data.AddUser(user)
	if err != nil {
		InvalidInputResponse(w, err.Error())
		return
	}

	output := templateRegisterSuccess{}
	output.Status = true
	output.Result.ID = int(id)
	output.Result.Username = inputJSONMap["username"].(string)

	outputJSON, err := json.Marshal(output)
	if err != nil {
		ServerError(w, err.Error())
		return
	}
	w.Write(outputJSON)
}
