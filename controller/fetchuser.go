package controller

import (
	"encoding/json"
	"loginsystem/data"
	"net/http"
	"strconv"
	"strings"
)

func FetchUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json; charset=utf-8")

	parsed := strings.Split(r.URL.Path, "/")
	if len(parsed) != 3 {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if parsed[2] == "" {
		InvalidInputResponse(w, "id not provided")
		return
	}

	userID, err := strconv.Atoi(parsed[2])
	if err != nil {
		InvalidInputResponse(w, err.Error())
		return
	}

	account, err := data.FetchUser(userID)
	if err != nil {
		InvalidInputResponse(w, err.Error())
		return
	}

	if len(account) == 0 {
		InvalidInputResponse(w, "user not found")
		return
	}

	output := templateUserData{}
	output.Status = true
	output.Result.ID = userID
	output.Result.Username = account["username"]
	output.Result.Name = account["name"]
	output.Result.Surname = account["surname"]

	outputJSON, err := json.Marshal(output)
	if err != nil {
		ServerErrorResponse(w, err.Error())
		return
	}
	w.Write(outputJSON)
}
