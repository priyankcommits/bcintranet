package controllers

import (
	"net/http"

	"bcintranet/templates"
	"bcintranet/utils"
)

func MetricsController(res http.ResponseWriter, req *http.Request) {
	// Show Metric graphs
	data := make(map[string]interface{})
	controllerTemplate := templates.METRICS
	if req.Method == "GET" {
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
}

func MetricsOpsAddController(res http.ResponseWriter, req *http.Request) {
	// Add and post a metric
	data := make(map[string]interface{})
	controllerTemplate := templates.METRICS
	if req.Method == "GET" {
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
	if req.Method == "POST" {

	}
}
