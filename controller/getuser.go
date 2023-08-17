package controller

import (
	"encoding/json"
	"loginsystem/data"
	"net/http"
	"strconv"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")

	userId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		InvalidResponse(w, "missing ID query")
		return
	}

	account := data.GetFromID(userId)

	if account == nil {
		InvalidResponse(w, "user not found")
		return
	}

	output := templateUserData{}
	output.Status = true
	output.Result.Id = account.Id
	output.Result.Username = account.Username
	output.Result.Name = account.Name
	output.Result.Surname = account.Surname

	outputJSON, _ := json.Marshal(output)
	w.Write(outputJSON)
}
