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

func MetricsController(res http.ResponseWriter, req *http.Request) {
	// Show Metric graphs
	data := make(map[string]interface{})
	controllerTemplate := templates.METRICS
	if req.Method == "GET" {
		dayMetrics, monthMetrics, yearMetrics, err := store.GetAttendanceMetrics()
		if err == nil {
			data["dayintime"] = dayMetrics.InTime
			data["dayouttime"] = dayMetrics.OutTime
			data["dayooo"] = dayMetrics.OOO
			data["daywfh"] = dayMetrics.WFH
			data["daytotal"] = dayMetrics.InTime + dayMetrics.OutTime + dayMetrics.OOO + dayMetrics.WFH
		}
		var count, inCount, outCount, oooCount, wfhCount int
		for _, item := range monthMetrics {
			count += 1
			inCount += item.InTime
			outCount += item.OutTime
			oooCount += item.OOO
			wfhCount += item.WFH
		}
		if inCount != 0 {
			data["monthintime"] = inCount / count
		} else {
			data["monthintime"] = 0
		}
		if outCount != 0 {
			data["monthouttime"] = outCount / count
		} else {
			data["monthouttime"] = 0
		}
		if oooCount != 0 {
			data["monthooo"] = oooCount / count
		} else {
			data["monthooo"] = 0
		}
		if wfhCount != 0 {
			data["monthwfh"] = wfhCount / count
		} else {
			data["monthwfh"] = 0
		}
		data["monthtotal"] = inCount + outCount + oooCount + wfhCount
		for _, item := range yearMetrics {
			count += 1
			inCount += item.InTime
			outCount += item.OutTime
			oooCount += item.OOO
			wfhCount += item.WFH
		}
		if inCount != 0 {
			data["yearintime"] = inCount / count
		} else {
			data["yearintime"] = 0
		}
		if outCount != 0 {
			data["yearouttime"] = outCount / count
		} else {
			data["yearouttime"] = 0
		}
		if oooCount != 0 {
			data["yearooo"] = oooCount / count
		} else {
			data["yearoo"] = 0
		}
		if wfhCount != 0 {
			data["yearwfh"] = wfhCount / count
		} else {
			data["yearwfh"] = 0
		}
		data["yeartotal"] = inCount + outCount + oooCount + wfhCount
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
				user, _ := store.GetUser(context.Get(req, "userid").(string))
				users, _ := store.GetAllUsers()
				for _, userTo := range users {
					ticker := new(models.MetricTicker)
					ticker.From = user
					ticker.To = userTo
					ticker.Type = 1
					ticker.Status = 0
					ticker.CreatedOn = time.Now()
					go store.CreateMetricTicker(ticker)
				}
				http.Redirect(res, req, urls.METRICS_PATH, http.StatusSeeOther)
			}
		}
	}
}
