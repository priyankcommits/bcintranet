package controllers

import (
	"net/http"
	"text/template"

	"bcintranet/templates"
)

func HomeController(res http.ResponseWriter, req *http.Request) {
	// Home/wall controller
	t, _ := template.ParseFiles(templates.BASE, templates.HOME)
	t.Execute(res, nil)
}
