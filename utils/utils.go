package utils

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func GetValidSession(req *http.Request) (*sessions.Session, error) {
	// Returns a valid authenticated user session
	sessStore := sessions.NewCookieStore([]byte("gplus"))
	return sessions.Store.Get(sessStore, req, "gplus_gothic_session")
}
