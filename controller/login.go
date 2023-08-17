package controller

import (
	"crypto/sha256"
	"encoding/json"
	"io"
	"loginsystem/data"
	"net/http"
	"reflect"
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

	account := data.GetFromName(inputJSONMap["username"].(string))
	if account == nil {
		InvalidResponse(w, "user not found")
		return
	}

	hashComputer := sha256.New()
	hashComputer.Write([]byte(inputJSONMap["password"].(string)))
	if !reflect.DeepEqual(account.Password, hashComputer.Sum(nil)) {
		InvalidResponse(w, "incorrect password")
		return
	}

	output := templateLoginSuccess{}
	output.Status = true
	output.Result.Id = account.Id
	output.Result.Username = account.Username

	outputJSON, _ := json.Marshal(output)
	w.Write(outputJSON)
}
