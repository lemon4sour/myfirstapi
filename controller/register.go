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

	inputJSON, _ := io.ReadAll(r.Body)

	inputJSONMap, err := JSONtoMap(inputJSON)
	if err != nil {
		InvalidResponse(w, err.Error())
		return
	}

	if !MatchesTemplate(inputJSONMap, templateRegisterRequest) {
		InvalidResponse(w, "invalid JSON parameters")
		return
	}

	id, err := data.Add(inputJSONMap)
	if err != nil {
		InvalidResponse(w, err.Error())
		return
	}

	output := templateRegisterSuccess{}
	output.Status = true
	output.Result.Id = id
	output.Result.Username = inputJSONMap["username"].(string)

	outputJSON, _ := json.Marshal(output)
	w.Write(outputJSON)
}
