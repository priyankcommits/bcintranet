package controllers

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/context"
)

func ProfileViewController(res http.ResponseWriter, req *http.Request) {
	// Profile View Controller
}

func ProfileEditController(res http.ResponseWriter, req *http.Request) {
	// Profile Edit Controller
	log.Println(context.Get(req, "userid"))
	if req.Method == "GET" {
		t, _ := template.ParseFiles("templates/base.html", "templates/profile.html")
		t.Execute(res, nil)
	}
	if req.Method == "POST" {
	}
}
