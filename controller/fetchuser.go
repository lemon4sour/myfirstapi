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

	token, exists := r.Header["Api-Token"]
	if !exists {
		Unauthorized(w)
		return
	}
	exists, err := data.SessionExists(token[0])
	if err != nil {
		InvalidInput(w, err.Error())
		return
	}
	if !exists {
		Unauthorized(w)
		return
	}

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
		InvalidInput(w, "id not provided")
		return
	}

	userID, err := strconv.Atoi(parsed[2])
	if err != nil {
		InvalidInput(w, err.Error())
		return
	}

	account, err := data.FetchUser(int64(userID))
	if err != nil {
		InvalidInput(w, err.Error())
		return
	}

	if len(account) == 0 {
		InvalidInput(w, "user not found")
		return
	}

	output := templateUserData{}
	output.Status = true
	output.Result.ID = int64(userID)
	output.Result.Username = account["username"]
	output.Result.Name = account["name"]
	output.Result.Surname = account["surname"]

	outputJSON, err := json.Marshal(output)
	if err != nil {
		ServerError(w, err.Error())
		return
	}
	w.Write(outputJSON)
}
