package controller

import (
	"encoding/json"
	"loginsystem/data"
	"net/http"
)

func GameResult(w http.ResponseWriter, r *http.Request) {
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

	if !MatchesTemplate(inputMap, templateGameResults) {
		InvalidInput(w, "invalid JSON parameters")
		return
	}

	uid1 := int64(inputMap["userid1"].(float64))
	uid2 := int64(inputMap["userid2"].(float64))
	score1 := inputMap["score1"].(float64)
	score2 := inputMap["score2"].(float64)

	ConcludeGame(uid1, uid2, score1, score2)

	output := templateScoreUpdateSuccess{}
	output.Status = true
	output.User1.ID = uid1
	user, err := data.FetchUser(uid1)
	if err != nil {
		ServerError(w, err.Error())
		return
	}
	output.User1.Username = user["username"]
	output.User1.Score = data.FetchScore(uid1)
	output.User2.ID = uid2
	user, err = data.FetchUser(uid2)
	if err != nil {
		ServerError(w, err.Error())
		return
	}
	output.User2.Username = user["username"]
	output.User2.Score = data.FetchScore(uid2)

	outputJSON, err := json.Marshal(output)
	if err != nil {
		ServerError(w, err.Error())
		return
	}
	w.Write(outputJSON)
}
