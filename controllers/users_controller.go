package controllers

import (
	"net/http"

	"bcintranet/store"
	"bcintranet/templates"
	"bcintranet/utils"
)

func UsersController(res http.ResponseWriter, req *http.Request) {
	// Users Controller
	data := make(map[string]interface{})
	controllerTemplate := templates.USERS
	if req.Method == "GET" {
		data["users"], _ = store.GetAllUsers()
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
}
