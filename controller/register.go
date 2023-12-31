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

	id, err := data.AddUser(inputJSONMap)
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
		ServerErrorResponse(w, err.Error())
		return
	}
	w.Write(outputJSON)
}
