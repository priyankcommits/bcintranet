package utils

import (
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"bcintranet/models"
	"bcintranet/store"
	"bcintranet/templates"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

func GetValidSession(req *http.Request) (*sessions.Session, error) {
	// Returns a valid authenticated user session
	sessStore := sessions.NewCookieStore([]byte(os.Getenv("bc_app_key")))
	return sessions.Store.Get(sessStore, req, "gplus_gothic_session")
}

func CustomTemplateExecute(res http.ResponseWriter, req *http.Request, templateName string, data map[string]interface{}) {
	// Append common templates and data structs and execute template
	t, _ := template.ParseFiles(templates.BASE, templates.NOTIFICATIONS, templates.TICKER, templateName)
	if len(data) == 0 {
		data = make(map[string]interface{})
		data["user"], _ = store.GetUser(context.Get(req, "userid").(string))
		data["isadmin"] = store.IsAdmin(context.Get(req, "userid").(string))
		data["notifications"], _ = store.GetNotifications(context.Get(req, "userid").(string))
		data["metrics"], _ = store.GetMetricTicker(context.Get(req, "userid").(string))
	} else {
		data["user"], _ = store.GetUser(context.Get(req, "userid").(string))
		data["isadmin"] = store.IsAdmin(context.Get(req, "userid").(string))
		data["notifications"], _ = store.GetNotifications(context.Get(req, "userid").(string))
		data["metrics"], _ = store.GetMetricTicker(context.Get(req, "userid").(string))
	}
	if err := t.Execute(res, data); err != nil {
		log.Println(err)
	}
}

func AddParamsToUrl(url string, args []models.Kwargs) string {
	// Add params to url using a splice of models.kwargs struct
	for _, arg := range args {
		url = strings.Replace(url, "{"+arg.Key+"}", arg.Value, 1)
	}
	return url
}
