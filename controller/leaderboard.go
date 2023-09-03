package controller

import (
	"encoding/json"
	"loginsystem/data"
	"net/http"
)

func Leaderboard(w http.ResponseWriter, r *http.Request) {
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

	if !MatchesTemplate(inputMap, templateLeaderboardRequest) {
		InvalidInput(w, "invalid JSON parameters")
		return
	}

	page := int(inputMap["page"].(float64))
	count := int(inputMap["count"].(float64))

	leaderboardList, err := data.LeaderboardPage(page, count)
	if err != nil {
		InvalidInput(w, err.Error())
		return
	}

	output := templateLeaderboardPage{}
	output.Status = true
	output.Result = make([]leaderboardPlacement, len(leaderboardList))
	for i, u := range leaderboardList {
		output.Result[i].ID = u.ID
		output.Result[i].Rank = u.Rank
		output.Result[i].Score = u.Score
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		ServerError(w, err.Error())
		return
	}
	w.Write(outputJSON)
}
