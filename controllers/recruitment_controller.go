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

func RecruitmentController(res http.ResponseWriter, req *http.Request) {
	// Recruitment metrics and list
	data := make(map[string]interface{})
	controllerTemplate := templates.RECRUITMENT
	if req.Method == "GET" {
		data["candidates"], _ = store.GetCandidates()
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
}

func CandidateAddController(res http.ResponseWriter, req *http.Request) {
	// Add a candidate
	data := make(map[string]interface{})
	controllerTemplate := templates.CANDIDATE_ADD
	if req.Method == "GET" {
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
	if req.Method == "POST" {
		err := req.ParseForm()
		candidate := new(models.Candidate)
		decoder := schema.NewDecoder()
		err = decoder.Decode(candidate, req.Form)
		if err != nil {
			data["message"] = models.Message{Value: "Something went wrong"}
			utils.CustomTemplateExecute(res, req, controllerTemplate, data)
		} else {
			err = store.CheckCandidate(candidate.Email)
			if err == nil {
				data["message"] = models.Message{Value: "Candidate with Email already exists"}
				utils.CustomTemplateExecute(res, req, controllerTemplate, data)
			} else {
				user, _ := store.GetUser(context.Get(req, "userid").(string))
				candidate.LastUpdateBy = user
				candidate.LastUpdateOn = time.Now()
				candidateObject, _ := store.AddCandidate(candidate)
				admins, _ := store.GetAdmins()
				for _, admin := range admins {
					if candidate.LastUpdateBy.UserID != admin.UserID {
						notification := new(models.Notification)
						notification.From = candidate.LastUpdateBy
						notification.To = admin
						notification.Type = 3
						notification.ResourceID = candidateObject.CandidateID.Hex()
						notification.Status = 0
						notification.CreatedOn = time.Now()
						go store.CreateNotification(notification)
						email := new(models.Email)
						email.To = admin
						email.Subject = user.FirstName + " added a new candidate"
						email.Body = "There has been a new candidate added to the recruitment portal.<a href='" + os.Getenv("bc_host") + "/admin/recruitment/candidate/" + candidateObject.CandidateID.Hex() + "/view/'>Click here to see</a><br><br>From, <br><br> BC Intranet"
						go helpers.SendGridEmail(email)
					}
				}
				http.Redirect(res, req, urls.RECRUITMENT_PATH, http.StatusSeeOther)
			}
		}
	}
}

func CandidateViewController(res http.ResponseWriter, req *http.Request) {
	// View a candidate
	data := make(map[string]interface{})
	controllerTemplate := templates.CANDIDATE_VIEW
	if req.Method == "GET" {
		var err error
		data["candidate"], err = store.GetCandidate(req.URL.Query().Get(":candidateid"))
		if err == nil {
			utils.CustomTemplateExecute(res, req, controllerTemplate, data)
		} else {
			http.Redirect(res, req, urls.RECRUITMENT_PATH, http.StatusSeeOther)
		}
	}
}

func CandidateEditController(res http.ResponseWriter, req *http.Request) {
	// Edit a candidate
	data := make(map[string]interface{})
	controllerTemplate := templates.CANDIDATE_EDIT
	if req.Method == "GET" {
		var err error
		data["candidate"], err = store.GetCandidate(req.URL.Query().Get(":candidateid"))
		if err == nil {
			utils.CustomTemplateExecute(res, req, controllerTemplate, data)
		}
	}
	if req.Method == "POST" {
		err := req.ParseForm()
		candidate, _ := store.GetCandidate(req.URL.Query().Get(":candidateid"))
		decoder := schema.NewDecoder()
		err = decoder.Decode(&candidate, req.Form)
		user, _ := store.GetUser(context.Get(req, "userid").(string))
		candidate.LastUpdateBy = user
		candidate.LastUpdateOn = time.Now()
		err = store.UpdateCandidate(req.URL.Query().Get(":candidateid"), candidate)
		if err == nil {
			var kwargs []models.Kwargs
			kwargs = append(kwargs, models.Kwargs{Key: "candidateid", Value: candidate.CandidateID.Hex()})
			http.Redirect(res, req, utils.AddParamsToUrl(urls.CANDIDATE_VIEW_PATH, kwargs), http.StatusSeeOther)
		} else {
			data["message"] = models.Message{Value: "Something went wrong"}
			utils.CustomTemplateExecute(res, req, controllerTemplate, data)
		}
	}
}

func CandidateCommentAddController(res http.ResponseWriter, req *http.Request) {
	// Add a comment on candidate profile
	if req.Method == "POST" {
		err := req.ParseForm()
		comment := new(models.CandidateComment)
		decoder := schema.NewDecoder()
		err = decoder.Decode(comment, req.Form)
		if err == nil {
			user, _ := store.GetUser(context.Get(req, "userid").(string))
			comment.User = user
			comment.Day = time.Now()
			err = store.SaveCandidateComment(comment)
		}
		var kwargs []models.Kwargs
		kwargs = append(kwargs, models.Kwargs{Key: "candidateid", Value: comment.CandidateID})
		http.Redirect(res, req, utils.AddParamsToUrl(urls.CANDIDATE_VIEW_PATH, kwargs), http.StatusSeeOther)
	}
}
