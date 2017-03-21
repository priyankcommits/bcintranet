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

	"github.com/gorilla/schema"
)

func MetricsController(res http.ResponseWriter, req *http.Request) {
	// Show Metric graphs
	data := make(map[string]interface{})
	controllerTemplate := templates.METRICS
	if req.Method == "GET" {
		data["dayintimecount"] = "10"
		data["dayouttimecount"] = "8"
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
}

func AdminMetricsOpsAddController(res http.ResponseWriter, req *http.Request) {
	// Add and post a metric
	data := make(map[string]interface{})
	controllerTemplate := templates.METRICS_OPS_ADD
	if req.Method == "GET" {
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
	if req.Method == "POST" {
		err := req.ParseForm()
		attendanceLog := new(models.MetricsAttendance)
		decoder := schema.NewDecoder()
		decoder.RegisterConverter(time.Time{}, helpers.ConvertFormDate)
		err = decoder.Decode(attendanceLog, req.Form)
		if err != nil {
			data["message"] = models.Message{Value: "Something went wrong, try again"}
			utils.CustomTemplateExecute(res, req, controllerTemplate, data)
		} else {
			if attendanceLog.Day.IsZero() {
				data["message"] = models.Message{Value: "Please enter a valid date"}
				utils.CustomTemplateExecute(res, req, controllerTemplate, data)
			} else {
				err = store.SaveAttendanceLog(attendanceLog)
				http.Redirect(res, req, urls.METRICS_PATH, http.StatusSeeOther)
			}
		}
	}
}
