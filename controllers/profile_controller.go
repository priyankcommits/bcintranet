package controllers

import (
	"net/http"

	"bcintranet/models"
	"bcintranet/store"
	"bcintranet/templates"
	"bcintranet/urls"
	"bcintranet/utils"

	"github.com/gorilla/context"
	"github.com/gorilla/schema"
)

func HomeController(res http.ResponseWriter, req *http.Request) {
	// Home/wall controller
	data := make(map[string]interface{})
	controllerTemplate := templates.HOME
	_, err := store.GetProfile(context.Get(req, "userid").(string))
	if err != nil {
		http.Redirect(res, req, urls.PROFILE_EDIT_PATH, http.StatusSeeOther)
	} else {
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
}

func ProfileViewController(res http.ResponseWriter, req *http.Request) {
	// Profile View Controller
	data := make(map[string]interface{})
	controllerTemplate := templates.PROFILE
	if req.Method == "GET" {
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
	if req.Method == "POST" {
	}
}

func ProfileEditController(res http.ResponseWriter, req *http.Request) {
	// Profile Edit Controller
	data := make(map[string]interface{})
	controllerTemplate := templates.PROFILE
	if req.Method == "GET" {
		utils.CustomTemplateExecute(res, req, controllerTemplate, nil)
	}
	if req.Method == "POST" {
		err := req.ParseForm()
		profile := new(models.Profile)
		decoder := schema.NewDecoder()
		err = decoder.Decode(profile, req.Form)
		if err != nil {
			data["message"] = models.Message{Value: "Something went wrong, try again"}
			utils.CustomTemplateExecute(res, req, controllerTemplate, data)
		} else {
			profile.UserID = context.Get(req, "userid").(string)
			err = store.SaveProfile(profile)
			if err != nil {
				data["message"] = models.Message{Value: "Something went wrong, try again"}
				utils.CustomTemplateExecute(res, req, controllerTemplate, data)
			} else {
				http.Redirect(res, req, urls.PROFILE_VIEW_PATH, http.StatusSeeOther)
			}
		}
	}
}
