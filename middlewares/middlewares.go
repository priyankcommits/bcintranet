package middlewares

import (
	// "log"
	"net/http"

	"bcintranet/store"
	"bcintranet/urls"
	"bcintranet/utils"

	"github.com/gorilla/context"
)

func GothLoginMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Retreiving session, redirecting if no session found
	session, _ := utils.GetValidSession(req)
	if session.Values["gplus"] == nil {
		http.Redirect(res, req, urls.ROOT_PATH, http.StatusTemporaryRedirect)
	}
	if session.Values["userid"] != "" {
		context.Set(req, "userid", session.Values["userid"])
	} else {
		http.Redirect(res, req, urls.LOGOUT_PATH, http.StatusTemporaryRedirect)
	}
	next(res, req)
}

func SetUserMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Appending the user id to every request and redirecting accordinly if no profile found
	session, _ := utils.GetValidSession(req)
	if session.Values["userid"] != "" {
		context.Set(req, "userid", session.Values["userid"])
	} else {
		http.Redirect(res, req, urls.LOGOUT_PATH, http.StatusTemporaryRedirect)
	}
	err := store.GetProfile(session.Values["userid"].(string))
	if err != nil {
		http.Redirect(res, req, urls.PROFILE_EDIT_PATH, http.StatusTemporaryRedirect)
	} else {
		http.Redirect(res, req, urls.HOME_PATH, http.StatusTemporaryRedirect)
	}
	next(res, req)
}
