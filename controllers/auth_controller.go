package controllers

import (
	// "log"
	"net/http"
	"text/template"

	"bcintranet/store"
	"bcintranet/templates"
	"bcintranet/urls"
	"bcintranet/utils"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func LoginController(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles(templates.LOGIN)
	t.Execute(res, nil)
}

func LogoutController(res http.ResponseWriter, req *http.Request) {
	session, _ := utils.GetValidSession(req)
	session.Options = &sessions.Options{Path: urls.ROOT_PATH, MaxAge: -1}
	session.Save(req, res)
	http.Redirect(res, req, urls.ROOT_PATH, http.StatusTemporaryRedirect)
}

func AuthController(res http.ResponseWriter, req *http.Request) {
	gothic.BeginAuthHandler(res, req)
}

func AuthCallbackController(res http.ResponseWriter, req *http.Request) {
	var gothUser goth.User
	gothUser, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		gothic.BeginAuthHandler(res, req)
	}
	session, _ := utils.GetValidSession(req)
	session.Values["userid"] = gothUser.UserID
	session.Save(req, res)
	_, err = store.GetUser(gothUser.UserID)
	if err != nil {
		store.CreateUserData(
			gothUser.UserID, gothUser.FirstName, gothUser.LastName,
			gothUser.Email, gothUser.AccessToken,
		)
	}
	http.Redirect(res, req, urls.HOME_PATH, http.StatusTemporaryRedirect)
}
