package controllers

import (
	"net/http"
	"text/template"

	"bcintranet/templates"
)

func NotFound(res http.ResponseWriter, req *http.Request) {
	// 404 controller
	t, _ := template.ParseFiles(templates.BASE, templates.NOTFOUND)
	t.Execute(res, nil)
}
