package controller

import (
	"loginsystem/data"
	"net/http"
)

func LogOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")

	token, verified := verifyToken(w, r)
	if !verified {
		return
	}

	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}

	if err := data.RemoveSession(token); err != nil {
		InvalidInput(w, err.Error())
		return
	}
}
