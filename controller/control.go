package controller

import (
	"encoding/json"
	"net/http"
	"reflect"
)

func InvalidResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	response := templateErrorResponse{Status: false, Message: "400 Bad Request - " + message}
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
