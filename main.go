package main

import (
	"loginsystem/controller"
	"net/http"
)

const PORT = ":8080"

func main() {
	sm := http.NewServeMux()
	sm.Handle("/register", http.HandlerFunc(controller.Register))
	sm.Handle("/login", http.HandlerFunc(controller.Login))
	sm.Handle("/user/", http.HandlerFunc(controller.GetUser))
	sm.Handle("/updateuser", http.HandlerFunc(controller.UpdateUser))
	http.ListenAndServe(PORT, sm)
}
