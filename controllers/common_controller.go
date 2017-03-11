package controllers

import (
	"net/http"

	"bcintranet/templates"
	"bcintranet/utils"
)

func NotFound(res http.ResponseWriter, req *http.Request) {
	// 404 controller
	utils.CustomTemplateExecute(res, req, templates.NOTFOUND, nil)
}
