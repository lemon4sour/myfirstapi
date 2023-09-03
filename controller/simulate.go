package controller

import (
	"encoding/json"
	"loginsystem/data"
	"math"
	"math/rand"
	"net/http"
)

func Simulate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")

	_, verified := verifyToken(w, r)
	if !verified {
		return
	}

	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}

	inputMap, exists := getMap(w, r)
	if !exists {
		return
	}

	if !MatchesTemplate(inputMap, templateSimulationParams) {
		InvalidInput(w, "invalid JSON parameters")
		return
	}

	userCount := int(inputMap["usercount"].(float64))
	userList := make([]int64, 0)
	for i := 0; i < userCount; i++ {
		id, err := data.GenerateUser()
		if err != nil {
			ServerError(w, err.Error())
			return
		}
		userList = append(userList, id)
	}

	for i := 0; i < userCount; i++ {
		for j := i + 1; j < userCount; j++ {
			score1 := math.Floor(rand.Float64() * 10)
			score2 := math.Floor(rand.Float64() * 10)
			ConcludeGame(userList[i], userList[j], score1, score2)
		}
	}

	output := templateSimulationResult{}
	output.Status = true
	output.Result = make([]userData, userCount)
	for i, userID := range userList {
		user, err := data.FetchUser(userID)
		if err != nil {
			ServerError(w, err.Error())
			return
		}
		output.Result[i].ID = int(userID)
		output.Result[i].Name = user["name"]
		output.Result[i].Surname = user["surname"]
		output.Result[i].Username = user["username"]
		output.Result[i].Score = data.FetchScore(userID)
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		ServerError(w, err.Error())
		return
	}
	w.Write(outputJSON)
}
