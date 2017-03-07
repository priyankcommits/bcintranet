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

func PayslipController(res http.ResponseWriter, req *http.Request) {
	// PaySlipController
	data := make(map[string]interface{})
	controllerTemplate := templates.PAYSLIP
	if req.Method == "GET" {
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
	if req.Method == "POST" {
		err := req.ParseForm()
		payslip := new(models.Payslip)
		decoder := schema.NewDecoder()
		decoder.RegisterConverter(time.Time{}, helpers.ConvertFormDate)
		err = decoder.Decode(payslip, req.Form)
		user, _ := store.GetUser(context.Get(req, "userid").(string))
		payslip.Requestor = user
		payslip.RequestedOn = time.Now()
		payslip.Status = 0
		err = store.SavePayslipRequest(payslip)
		admins, _ := store.GetAdmins()
		for _, admin := range admins {
			if admin.UserID != user.UserID {
				notification := new(models.Notification)
				notification.From = user
				notification.To = admin
				notification.Type = 4
				notification.ResourceID = user.UserID
				notification.CreatedOn = time.Now()
				go store.CreateNotification(notification)
				email := new(models.Email)
				email.To = admin
				email.Subject = user.FirstName + " requested a payslip approval"
				email.Body = user.FirstName + " has requested a payslip. <a href='" + os.Getenv("bc_host") + "/admin/approvals/'>Click here to approve</a><br><br>From, <br><br> BC Inttranet"
				go helpers.SendGridEmail(email)
			}
		}
		if err == nil {
			http.Redirect(res, req, urls.PAYSLIP_HISTORY_PATH, http.StatusSeeOther)
		} else {
			utils.CustomTemplateExecute(res, req, controllerTemplate, data)
		}
	}
}

func PayslipHistoryController(res http.ResponseWriter, req *http.Request) {
	// PaySlip History Controller
	data := make(map[string]interface{})
	controllerTemplate := templates.PAYSLIP_HISTORY
	if req.Method == "GET" {
		payslips, _ := store.GetPayslipHistory(context.Get(req, "userid").(string))
		for _, payslip := range payslips {
			if payslip.Status == 1 {
				helpers.GeneratePayslipPDF(payslip)
			}
		}
		data["payslips"] = payslips
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
}
