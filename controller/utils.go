package controller

import (
	"encoding/json"
	"log"
	"loginsystem/data"
	"net/http"
	"os"
	"reflect"
)

func init() {
	log.SetOutput(os.Stdout)
}

func InvalidInputResponse(w http.ResponseWriter, message string) {
	log.Println(message)
	w.WriteHeader(http.StatusBadRequest)
	response := templateErrorResponse{Status: false, Message: "400 Bad Request - " + message}
	mresponse, _ := json.Marshal(response)
	w.Write(mresponse)
}

func ServerError(w http.ResponseWriter, message string) {
	log.Println(message)
	w.WriteHeader(http.StatusInternalServerError)
	response := templateErrorResponse{Status: false, Message: "500 Internal Server Error - " + message}
	mresponse, _ := json.Marshal(response)
	w.Write(mresponse)
}

func JSONtoMap(data []byte) (map[string]any, error) {
	var JSONMap map[string]any
	err := json.Unmarshal(data, &JSONMap)
	if err != nil {
		return nil, err
	}
	return JSONMap, nil
}

func MatchesTemplate(data map[string]any, template map[string]any) bool {
	/*
		if len(data) != len(template) {
			return false
		}
	*/
	for key, tvalue := range template {
		value, exists := data[key]
		if !exists {
			return false
		}
		if reflect.TypeOf(value) != reflect.TypeOf(tvalue) {
			return false
		}
	}
	return true
}

func ConcludeGame(userid1, userid2 int64, score1, score2 float64) {
	if score1 > score2 {
		data.AddScore(userid1, 3)
	} else if score1 < score2 {
		data.AddScore(userid2, 3)
	} else {
		data.AddScore(userid1, 1)
		data.AddScore(userid2, 1)
	}
}
