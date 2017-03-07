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

func HomeController(res http.ResponseWriter, req *http.Request) {
	// Home/wall controller, all posts
	data := make(map[string]interface{})
	controllerTemplate := templates.HOME
	if req.Method == "GET" {
		_, err := store.GetProfile(context.Get(req, "userid").(string))
		if err != nil {
			http.Redirect(res, req, urls.PROFILE_EDIT_PATH, http.StatusSeeOther)
		} else {
			data["posts"], _ = store.GetAllPosts()
			utils.CustomTemplateExecute(res, req, controllerTemplate, data)
		}
	}
}

func AnnouncementsController(res http.ResponseWriter, req *http.Request) {
	// Announcements Only controller, announcemnet posts only
	data := make(map[string]interface{})
	controllerTemplate := templates.ANNOUNCEMENTS
	if req.Method == "GET" {
		_, err := store.GetProfile(context.Get(req, "userid").(string))
		if err != nil {
			http.Redirect(res, req, urls.PROFILE_EDIT_PATH, http.StatusSeeOther)
		} else {
			data["posts"], _ = store.GetAnnouncements()
			utils.CustomTemplateExecute(res, req, controllerTemplate, data)
		}
	}
}

func WallPostController(res http.ResponseWriter, req *http.Request) {
	// Show one post
	data := make(map[string]interface{})
	controllerTemplate := templates.WALL_POST
	if req.Method == "GET" {
		var err error
		data["post"], err = store.GetPost(req.URL.Query().Get(":postid"))
		if err != nil {
			queryParam := "?m=Unable to find resource"
			http.Redirect(res, req, urls.HOME_PATH+queryParam, http.StatusSeeOther)
		} else {
			utils.CustomTemplateExecute(res, req, controllerTemplate, data)
		}
	}
}

func WallPostAddController(res http.ResponseWriter, req *http.Request) {
	// Add a post
	data := make(map[string]interface{})
	controllerTemplate := templates.POST_ADD
	if req.Method == "GET" {
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
	if req.Method == "POST" {
		err := req.ParseForm()
		post := new(models.WallPost)
		decoder := schema.NewDecoder()
		err = decoder.Decode(post, req.Form)
		if err != nil {
			data["message"] = models.Message{Value: "Something went wrong, try again"}
			utils.CustomTemplateExecute(res, req, controllerTemplate, data)
		} else {
			user, _ := store.GetUser(context.Get(req, "userid").(string))
			post.User = user
			post.Day = time.Now()
			post, err = store.SavePost(post)
			if post.Type == "2" {
				users, _ := store.GetAllUsers()
				for _, userTo := range users {
					if post.User.UserID != userTo.UserID {
						notification := new(models.Notification)
						notification.From = post.User
						notification.To = userTo
						notification.Type = 2
						notification.ResourceID = post.PostID.Hex()
						notification.Status = 0
						notification.CreatedOn = time.Now()
						go store.CreateNotification(notification)
						email := new(models.Email)
						email.To = userTo
						email.Subject = userTo.FirstName + " made an announcement"
						email.Body = "There's been a new announcement <a href='" + os.Getenv("bc_host") + "/bc/post/" + post.PostID.Hex() + "/view/'>Click here to see</a><br><br>From, <br><br>BC Intranet'"
						go helpers.SendGridEmail(email)
					}
				}
			}
			var kwargs []models.Kwargs
			kwargs = append(kwargs, models.Kwargs{Key: "postid", Value: post.PostID.Hex()})
			http.Redirect(res, req, utils.AddParamsToUrl(urls.WALL_POST_PATH, kwargs), http.StatusSeeOther)
		}
	}
}

func WallPostDeleteController(res http.ResponseWriter, req *http.Request) {
	// Delete a post
	if req.Method == "GET" {
		post, err := store.GetPost(req.URL.Query().Get(":postid"))
		if err == nil {
			if post.User.UserID == context.Get(req, "userid").(string) {
				store.DeletePost(post.PostID.Hex())
				http.Redirect(res, req, urls.HOME_PATH, http.StatusSeeOther)
			} else {
				http.Redirect(res, req, urls.HOME_PATH, http.StatusSeeOther)
			}
		} else {
			http.Redirect(res, req, urls.HOME_PATH, http.StatusSeeOther)
		}
	}
}

func WallPostCommentAddController(res http.ResponseWriter, req *http.Request) {
	// Add comment to post
	if req.Method == "POST" {
		err := req.ParseForm()
		comment := new(models.WallPostComment)
		decoder := schema.NewDecoder()
		err = decoder.Decode(comment, req.Form)
		if err == nil {
			user, _ := store.GetUser(context.Get(req, "userid").(string))
			post, _ := store.GetPost(comment.PostID)
			comment.User = user
			comment.Day = time.Now()
			err = store.SaveComment(comment)
			if post.User.UserID != user.UserID {
				notification := new(models.Notification)
				notification.From = user
				notification.To = post.User
				notification.Type = 1
				notification.ResourceID = post.PostID.Hex()
				notification.Status = 0
				notification.CreatedOn = time.Now()
				go store.CreateNotification(notification)
				email := new(models.Email)
				email.To = post.User
				email.Subject = user.FirstName + " commented on you post"
				email.Body = "You have a new comment on your post. <a href='" + os.Getenv("bc_host") + "/bc/post/" + post.PostID.Hex() + "/view/'>Click here to see</a> <br><br> From, <br><br> BC Intranet"
				go helpers.SendGridEmail(email)
			}
			var kwargs []models.Kwargs
			kwargs = append(kwargs, models.Kwargs{Key: "postid", Value: comment.PostID})
			http.Redirect(res, req, utils.AddParamsToUrl(urls.WALL_POST_PATH, kwargs), http.StatusSeeOther)
		}
	}
}

func WallPostCommentDeleteController(res http.ResponseWriter, req *http.Request) {
	// Delete a comment
	if req.Method == "GET" {
		post, err := store.GetPost(req.URL.Query().Get(":postid"))
		if err == nil {
			for _, item := range post.Comments {
				if item.User.UserID == context.Get(req, "userid") && item.CommentID.Hex() == req.URL.Query().Get(":commentid") {
					store.DeleteComment(post.PostID.Hex(), req.URL.Query().Get(":commentid"))
				}
			}
			http.Redirect(res, req, urls.HOME_PATH, http.StatusSeeOther)
		} else {
			http.Redirect(res, req, urls.HOME_PATH, http.StatusSeeOther)
		}
	}
}
