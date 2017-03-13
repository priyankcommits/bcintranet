package controllers

import (
	"net/http"

	"bcintranet/templates"
	"bcintranet/utils"
)

func UsersController(res http.ResponseWriter, req *http.Request) {
	// Users Controller
	data := make(map[string]interface{})
	controllerTemplate := templates.USERS
	utils.CustomTemplateExecute(res, req, controllerTemplate, data)
}
