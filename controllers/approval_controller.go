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

func AdminApprovalController(res http.ResponseWriter, req *http.Request) {
	// Approvals landing page Controller
	data := make(map[string]interface{})
	controllerTemplate := templates.APPROVALS
	if req.Method == "GET" {
		data["payslips"], _ = store.GetAllPayslips()
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
}

func AdminApprovalPayslipController(res http.ResponseWriter, req *http.Request) {
	// Payslip approval Controller
	data := make(map[string]interface{})
	controllerTemplate := templates.APPROVALS_PAYSLIP
	if req.Method == "GET" {
		data["payslip"], _ = store.GetPayslip(req.URL.Query().Get(":payslipid"))
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
	if req.Method == "POST" {
		err := req.ParseForm()
		payslip, _ := store.GetPayslip(req.URL.Query().Get(":payslipid"))
		decoder := schema.NewDecoder()
		err = decoder.Decode(&payslip, req.Form)
		user, _ := store.GetUser(context.Get(req, "userid").(string))
		payslip.Approver = user
		payslip.Status = 1
		err = store.UpdatePayslip(req.URL.Query().Get(":payslipid"), payslip)
		if err != nil {
			data["message"] = models.Message{Value: "Something went wrong"}
			utils.CustomTemplateExecute(res, req, controllerTemplate, data)
		} else {
			http.Redirect(res, req, urls.APPROVALS_PATH, http.StatusSeeOther)
		}
	}
}
