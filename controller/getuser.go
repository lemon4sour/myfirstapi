package controller

import (
	"encoding/json"
	"loginsystem/data"
	"net/http"
	"strconv"
	"strings"
)

func GetUser(w http.ResponseWriter, r *http.Request) {

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
		InvalidResponse(w, "id not provided")
		return
	}

	userId, err := strconv.Atoi(parsed[2])
	if err != nil {
		InvalidResponse(w, err.Error())
		return
	}

	account := data.GetUser(userId)

	if len(account) == 0 {
		InvalidResponse(w, "user not found")
		return
	}

	output := templateUserData{}
	output.Status = true
	output.Result.Id = userId
	output.Result.Username = account["username"]
	output.Result.Name = account["name"]
	output.Result.Surname = account["surname"]

	outputJSON, _ := json.Marshal(output)
	w.Write(outputJSON)
}
