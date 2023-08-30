package controller

import (
	"encoding/json"
	"io"
	"loginsystem/data"
	"math/rand"
	"net/http"
)

func Simulate(w http.ResponseWriter, r *http.Request) {
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

	if !MatchesTemplate(inputJSONMap, templateSimulationParams) {
		InvalidInputResponse(w, "invalid JSON parameters")
		return
	}

	userCount := int(inputJSONMap["usercount"].(float64))
	userList := make([]int64, 0)
	for i := 0; i < userCount; i++ {
		id, err := data.GenerateUser()
		if err != nil {
			ServerErrorResponse(w, err.Error())
			return
		}
		userList = append(userList, id)
	}
	for i := 0; i < userCount; i++ {
		for j := i + 1; j < userCount; j++ {
			score1 := rand.Float64() * 10
			score2 := rand.Float64() * 10

			ConcludeGame(int(userList[i]), int(userList[j]), score1, score2)
		}
	}

	output := templateSimulationResult{}
	output.Status = true
	output.Result = make([]userData, userCount)
	for i, userID := range userList {
		user, err := data.FetchUser(int(userID))
		if err != nil {
			ServerErrorResponse(w, err.Error())
			return
		}
		output.Result[i].ID = int(userID)
		output.Result[i].Name = user["name"]
		output.Result[i].Surname = user["surname"]
		output.Result[i].Username = user["username"]
		output.Result[i].Score = data.FetchScore(int(userID))
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		ServerErrorResponse(w, err.Error())
		return
	}
	w.Write(outputJSON)
}
