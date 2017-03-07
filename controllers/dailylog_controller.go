package controllers

import (
	"net/http"
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

func DailyLogsController(res http.ResponseWriter, req *http.Request) {
	// Show Metric graphs
	data := make(map[string]interface{})
	controllerTemplate := templates.DAILY_LOGS
	if req.Method == "GET" {
		dayLogCount, monthLogCount, yearLogCount, _ := store.GetDailyLogMetrics()
		data["daycount"] = dayLogCount
		data["monthcount"] = monthLogCount
		data["yearcount"] = yearLogCount
		data["users"], _ = store.GetAllUsers()
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
}

func DailyLogsAddController(res http.ResponseWriter, req *http.Request) {
	// Add a log entry
	data := make(map[string]interface{})
	controllerTemplate := templates.DAILY_LOGS_ADD
	if req.Method == "GET" {
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
	if req.Method == "POST" {
		err := req.ParseForm()
		dailyLog := new(models.MetricsDailyLogs)
		dailyLog.UserID = context.Get(req, "userid").(string)
		decoder := schema.NewDecoder()
		decoder.RegisterConverter(time.Time{}, helpers.ConvertFormDate)
		err = decoder.Decode(dailyLog, req.Form)
		if err != nil {
			data["message"] = models.Message{Value: "Something went wrong, try again"}
			utils.CustomTemplateExecute(res, req, controllerTemplate, data)
		} else {
			if dailyLog.Day.IsZero() {
				data["message"] = models.Message{Value: "Please enter a valid date"}
				utils.CustomTemplateExecute(res, req, controllerTemplate, data)
			} else {
				err = store.SaveDailyLog(dailyLog)
				user, _ := store.GetUser(context.Get(req, "userid").(string))
				users, _ := store.GetAllUsers()
				for _, userTo := range users {
					ticker := new(models.MetricTicker)
					ticker.From = user
					ticker.To = userTo
					ticker.ResourceID = user.UserID
					ticker.Type = 2
					ticker.Status = 0
					ticker.CreatedOn = time.Now()
					go store.CreateMetricTicker(ticker)
				}
				http.Redirect(res, req, urls.DAILY_LOGS_PATH, http.StatusSeeOther)
			}
		}
	}
}

func DailyLogsViewController(res http.ResponseWriter, req *http.Request) {
	// Display All Entries made by the user
	data := make(map[string]interface{})
	controllerTemplate := templates.DAILY_LOGS_VIEW
	if req.Method == "GET" {
		data["requested_user"], _ = store.GetUser(req.URL.Query().Get(":userid"))
		data["entries"], _ = store.GetUserDailyLogs(req.URL.Query().Get(":userid"))
		data["users"], _ = store.GetAllUsers()
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
}
