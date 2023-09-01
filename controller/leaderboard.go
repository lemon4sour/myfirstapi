package controller

import (
	"encoding/json"
	"io"
	"loginsystem/data"
	"net/http"
)

func Leaderboard(w http.ResponseWriter, r *http.Request) {
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

	if !MatchesTemplate(inputJSONMap, templateLeaderboardRequest) {
		InvalidInputResponse(w, "invalid JSON parameters")
		return
	}

	page := int(inputJSONMap["page"].(float64))
	count := int(inputJSONMap["count"].(float64))

	leaderboardList, err := data.LeaderboardPage(page, count)
	if err != nil {
		InvalidInputResponse(w, err.Error())
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
