package controllers

import (
	"net/http"
	"strings"
	"text/template"

	"bcintranet/store"
	"bcintranet/templates"
	"bcintranet/urls"
	"bcintranet/utils"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func LoginController(res http.ResponseWriter, req *http.Request) {
	// login page controller
	session, _ := utils.GetValidSession(req)
	if session.Values["userid"] != nil {
		context.Set(req, "userid", session.Values["userid"])
		http.Redirect(res, req, urls.HOME_PATH, http.StatusSeeOther)
	} else {
		t, _ := template.ParseFiles(templates.LOGIN)
		t.Execute(res, nil)
	}
}

func LogoutController(res http.ResponseWriter, req *http.Request) {
	// delete the cookie and redirect
	session, _ := utils.GetValidSession(req)
	session.Options = &sessions.Options{Path: urls.ROOT_PATH, MaxAge: -1}
	session.Save(req, res)
	http.Redirect(res, req, urls.ROOT_PATH, http.StatusTemporaryRedirect)
}

func AuthController(res http.ResponseWriter, req *http.Request) {
	// start authentication process using gothic
	gothic.BeginAuthHandler(res, req)
}

func AuthCallbackController(res http.ResponseWriter, req *http.Request) {
	// goth callback controller to complete user auth and create user
	var gothUser goth.User
	gothUser, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		gothic.BeginAuthHandler(res, req)
	}
	emailDomain := strings.Join(strings.Split(gothUser.Email, "@")[1:], "a")
	if emailDomain == "" {
		session, _ := utils.GetValidSession(req)
		session.Options = &sessions.Options{Path: urls.ROOT_PATH, MaxAge: -1}
		session.Save(req, res)
		queryParam := "?m=account_invalid"
		http.Redirect(res, req, urls.ROOT_PATH+queryParam, http.StatusSeeOther)
	}
	session, _ := utils.GetValidSession(req)
	session.Values["userid"] = gothUser.UserID
	session.Save(req, res)
	userAdmin := false
	adminList := []string{
		"pulumati.priyank@gmail.com",
	}
	for _, value := range adminList {
		if value == gothUser.Email {
			userAdmin = true
		}
	}
	store.SaveUser(
		gothUser.UserID, gothUser.FirstName, gothUser.LastName,
		gothUser.Email, gothUser.AccessToken, gothUser.AvatarURL,
		userAdmin,
	)
	http.Redirect(res, req, urls.HOME_PATH, http.StatusSeeOther)
}
