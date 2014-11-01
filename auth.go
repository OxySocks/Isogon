package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/render"
	"net/http"
)

// Martini handler that handles session validation.
func sessionValidator() martini.Handler {
	return func(session sessions.Session) {
		data := session.Get("user")
		_, exists := data.(int64)
		if exists {
			return
		}
		session.Set("user", -1)
		return
	}
}

// Function that checks if the current session is available and if it can be validated.
func sessionIsValid(session sessions.Session) bool {
	data := session.Get("user")
	_, exists := data.(int64)
	if exists {
		return true
	}
	return false
}

// ProtectedPage is a martini handler that makes sure users can only enter a page authorized.
// Should be included as a parameter for routes that should be secured.
func ProtectedPage(req *http.Request, session sessions.Session, r render.Render) {
	if !sessionIsValid(session) {
		session.Delete("user")
		r.Error(401)
	}
}
