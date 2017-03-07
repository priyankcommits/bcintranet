package controllers

import (
	"log"
	"net/http"
	"text/template"

	"bcintranet/models"
	"bcintranet/store"
	"bcintranet/templates"
	"bcintranet/urls"

	"github.com/gorilla/context"
	"github.com/gorilla/schema"
)

func ProfileViewController(res http.ResponseWriter, req *http.Request) {
	// Profile View Controller
	if req.Method == "GET" {
		log.Println("in profile get view")
	}
	if req.Method == "POST" {
	}
}

func ProfileEditController(res http.ResponseWriter, req *http.Request) {
	// Profile Edit Controller
	t, _ := template.ParseFiles(templates.BASE, templates.PROFILE)
	data := make(map[string]interface{})
	data["user"], _ = store.GetUser(context.Get(req, "userid").(string))
	if req.Method == "GET" {
		t.Execute(res, data)
	}
	if req.Method == "POST" {
		err := req.ParseForm()
		profile := new(models.Profile)
		decoder := schema.NewDecoder()
		err = decoder.Decode(profile, req.Form)
		if err != nil {
			data["message"] = models.Message{Value: "Something went wrong, try again"}
			t.Execute(res, data)
		} else {
			profile.UserID = context.Get(req, "userid").(string)
			err = store.CreateProfile(profile)
			if err != nil {
				data["message"] = models.Message{Value: "Something went wrong, try again"}
				t.Execute(res, data)
			} else {
				http.Redirect(res, req, urls.PROFILE_VIEW_PATH, http.StatusSeeOther)
			}
		}
	}
}
