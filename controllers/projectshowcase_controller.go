package controllers

import (
	"net/http"
	//"time"

	//"bcintranet/helpers"
	//"bcintranet/models"
	//"bcintranet/store"
	"bcintranet/templates"
	//"bcintranet/urls"
	"bcintranet/utils"

)

func ProjectShowcaseController(res http.ResponseWriter, req *http.Request) {
	// Projects Views
	data := make(map[string]interface{})
	controllerTemplate := templates.PROJECT_SHOWCASE_VIEW
	if req.Method == "GET" {
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
}
