package controller

import (
	"encoding/json"
	"io"
	"log"
	"loginsystem/data"
	"net/http"
	"os"
	"reflect"
)

var apiTokenKey = "Api-Token"

func init() {
	log.SetOutput(os.Stdout)
}

func InvalidInput(w http.ResponseWriter, message string) {
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

func Unauthorized(w http.ResponseWriter) {
	log.Println("invalid or missing token")
	w.WriteHeader(http.StatusUnauthorized)
	response := templateErrorResponse{Status: false, Message: "401 Unauthorized - invalid or missing token"}
	mresponse, _ := json.Marshal(response)
	w.Write(mresponse)
}

func methodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
}

func getMap(w http.ResponseWriter, r *http.Request) (map[string]any, bool) {
	inputJSON, err := io.ReadAll(r.Body)
	if err != nil {
		InvalidInput(w, err.Error())
		return nil, false
	}

	inputJSONMap, err := JSONtoMap(inputJSON)
	if err != nil {
		InvalidInput(w, err.Error())
		return nil, false
	}

	return inputJSONMap, true
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
		data.Log(userid1, "won")
		data.Log(userid2, "lost")
		data.AddScore(userid1, 3)
	} else if score1 < score2 {
		data.Log(userid1, "lost")
		data.Log(userid2, "won")
		data.AddScore(userid2, 3)
	} else {
		data.Log(userid1, "tie")
		data.Log(userid2, "tie")
		data.AddScore(userid1, 1)
		data.AddScore(userid2, 1)
	}
}

func verifyToken(w http.ResponseWriter, r *http.Request) (string, bool) {
	token, exists := r.Header["Api-Token"]
	if !exists {
		Unauthorized(w)
		return "", false
	}
	exists, err := data.SessionExists(token[0])
	if err != nil {
		InvalidInput(w, err.Error())
		return "", false
	}
	if !exists {
		Unauthorized(w)
		return "", false
	}
	return token[0], true
}
