package main

import (
	"loginsystem/controller"
	"net/http"
)

const port = ":8080"

func main() {
	sm := http.NewServeMux()
	sm.Handle("/register", http.HandlerFunc(controller.Register))
	sm.Handle("/login", http.HandlerFunc(controller.Login))
	sm.Handle("/user/", http.HandlerFunc(controller.FetchUser))
	sm.Handle("/rename", http.HandlerFunc(controller.Rename))
	sm.Handle("/gameresult", http.HandlerFunc(controller.GameResult))
	sm.Handle("/leaderboard", http.HandlerFunc(controller.Leaderboard))
	sm.Handle("/simulate", http.HandlerFunc(controller.Simulate))
	http.ListenAndServe(port, sm)
}
