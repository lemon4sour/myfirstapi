package controller

import (
	"encoding/json"
	"io"
	"loginsystem/data"
	"net/http"
)

func GameResult(w http.ResponseWriter, r *http.Request) {
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

	if !MatchesTemplate(inputJSONMap, templateGameResults) {
		InvalidInputResponse(w, "invalid JSON parameters")
		return
	}

	uid1 := int(inputJSONMap["userid1"].(float64))
	uid2 := int(inputJSONMap["userid2"].(float64))
	score1 := inputJSONMap["score1"].(float64)
	score2 := inputJSONMap["score2"].(float64)

	ConcludeGame(uid1, uid2, score1, score2)

	output := templateScoreUpdateSuccess{}
	output.Status = true
	output.User1.ID = uid1
	user, err := data.FetchUser(uid1)
	if err != nil {
		ServerErrorResponse(w, err.Error())
		return
	}
	output.User1.Username = user["username"]
	output.User1.Score = data.FetchScore(uid1)
	output.User2.ID = uid2
	user, err = data.FetchUser(uid2)
	if err != nil {
		ServerErrorResponse(w, err.Error())
		return
	}
	output.User2.Username = user["username"]
	output.User2.Score = data.FetchScore(uid2)

	outputJSON, err := json.Marshal(output)
	if err != nil {
		ServerErrorResponse(w, err.Error())
		return
	}
	w.Write(outputJSON)
}
