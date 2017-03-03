package utils

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func GetValidSession(req *http.Request) (*sessions.Session, error) {
	sess_store := sessions.NewCookieStore([]byte("gplus"))
	return sessions.Store.Get(sess_store, req, "gplus_gothic_session")
}
