package controllers

import (
	"net/http"
	"text/template"
)

func NotFound(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/base.html", "templates/404.html")
	t.Execute(res, nil)
}
