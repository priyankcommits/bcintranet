package controllers

import (
	"net/http"
	"os"
	"time"

	"bcintranet/helpers"
	"bcintranet/models"
	"bcintranet/store"
	"bcintranet/templates"
	"bcintranet/urls"
	"bcintranet/utils"

	"github.com/gorilla/context"
	"github.com/gorilla/schema"
)

func AboutController(res http.ResponseWriter, req *http.Request) {
	// About site controller
	data := make(map[string]interface{})
	controllerTemplate := templates.ABOUT
	if req.Method == "GET" {
		data["about"], _ = store.GetAbout()
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
}

func AboutFeedbackController(res http.ResponseWriter, req *http.Request) {
	// Get feedback
	if req.Method == "POST" {
		err := req.ParseForm()
		if err == nil {
			user, err := store.GetUser("101507997126494782640")
			if err == nil {
				email := new(models.Email)
				email.To = user
				email.Subject = user.FirstName + " has some feedback for BC Intranet"
				email.Body = "Type: " + req.Form["Feedback"][0] + "<br><br>" + "Message: " + req.Form["Message"][0]
				go helpers.SendGridEmail(email)
			}
		}
		queryParam := "?m=Feedback Sent"
		http.Redirect(res, req, urls.ABOUT_PATH+queryParam, http.StatusSeeOther)
	}
}

func ProjectsController(res http.ResponseWriter, req *http.Request) {
	// Show all projects
	data := make(map[string]interface{})
	controllerTemplate := templates.PROJECTS
	if req.Method == "GET" {
		data["projects"], _ = store.GetProjects()
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
}

func ProjectAddController(res http.ResponseWriter, req *http.Request) {
	// Project add
	data := make(map[string]interface{})
	controllerTemplate := templates.PROJECT_ADD
	if req.Method == "GET" {
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
	if req.Method == "POST" {
		err := req.ParseForm()
		project := new(models.Project)
		decoder := schema.NewDecoder()
		err = decoder.Decode(project, req.Form)
		if err != nil {
			http.Redirect(res, req, urls.PROJECTS_PATH, http.StatusSeeOther)
		} else {
			user, _ := store.GetUser(context.Get(req, "userid").(string))
			project.User = user
			project, err = store.SaveProject(project)
			users, _ := store.GetAllUsers()
			for _, userTo := range users {
				if user.UserID != userTo.UserID {
					ticker := new(models.MetricTicker)
					ticker.From = user
					ticker.To = userTo
					ticker.ResourceID = project.ProjectID.Hex()
					ticker.Type = 3
					ticker.CreatedOn = time.Now()
					go store.CreateMetricTicker(ticker)
					email := new(models.Email)
					email.To = userTo
					email.Subject = user.FirstName + " added a new project"
					email.Body = "A new project has been added !! <a href='" + os.Getenv("bc_host") + "/bc/project/" + project.ProjectID.Hex() + "/view/'>Click here to see</a> <br><br> From, <br><br> BC Intranet"
					go helpers.SendGridEmail(email)
				}
			}
			http.Redirect(res, req, urls.PROJECTS_PATH, http.StatusSeeOther)
		}
	}
}

func ProjectViewController(res http.ResponseWriter, req *http.Request) {
	// View a Project
	data := make(map[string]interface{})
	controllerTemplate := templates.PROJECT_VIEW
	if req.Method == "GET" {
		data["project"], _ = store.GetProject(req.URL.Query().Get(":projectid"))
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
}

func ProjectEditController(res http.ResponseWriter, req *http.Request) {
	// Edit a project
	data := make(map[string]interface{})
	controllerTemplate := templates.PROJECT_EDIT
	if req.Method == "GET" {
		data["project"], _ = store.GetProject(req.URL.Query().Get(":projectid"))
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
	if req.Method == "POST" {
		err := req.ParseForm()
		project, _ := store.GetProject(req.URL.Query().Get(":projectid"))
		if err == nil {
			user, _ := store.GetUser(context.Get(req, "userid").(string))
			if project.User.UserID == user.UserID {
				decoder := schema.NewDecoder()
				err = decoder.Decode(&project, req.Form)
				err = store.UpdateProject(project.ProjectID.Hex(), project)
				if err != nil {
					data["message"] = models.Message{Value: "Something went wrong"}
					utils.CustomTemplateExecute(res, req, controllerTemplate, data)
				}
			}
		}
		queryParam := "?m=Updated Project"
		var kwargs []models.Kwargs
		kwargs = append(kwargs, models.Kwargs{Key: "projectid", Value: project.ProjectID.Hex()})
		http.Redirect(res, req, utils.AddParamsToUrl(urls.PROJECT_VIEW_PATH, kwargs)+queryParam, http.StatusSeeOther)
	}
}

func ProjectFeedbackController(res http.ResponseWriter, req *http.Request) {
	// Get project feedback
	if req.Method == "POST" {
		err := req.ParseForm()
		project, _ := store.GetProject(req.URL.Query().Get(":projectid"))
		if err == nil {
			user, err := store.GetUser(context.Get(req, "userid").(string))
			if err == nil {
				email := new(models.Email)
				email.To = project.User
				email.Subject = user.FirstName + " has some feedback for" + project.Title
				email.Body = "Type: " + req.Form["Feedback"][0] + "<br><br>" + "Message: " + req.Form["Message"][0]
				go helpers.SendGridEmail(email)
			}
		}
		var kwargs []models.Kwargs
		kwargs = append(kwargs, models.Kwargs{Key: "projectid", Value: project.ProjectID.Hex()})
		queryParam := "?m=Feedback Sent"
		http.Redirect(res, req, utils.AddParamsToUrl(urls.PROJECT_VIEW_PATH, kwargs)+queryParam, http.StatusSeeOther)
	}
}
