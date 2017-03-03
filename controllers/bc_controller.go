package controllers

import (
	"net/http"
	"text/template"
	// "github.com/gorilla/context"
)

func HomeController(res http.ResponseWriter, req *http.Request) {
	// Home/wall controller
	t, _ := template.ParseFiles("templates/base.html", "templates/home.html")
	t.Execute(res, nil)
}
