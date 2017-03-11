package utils

import (
	"net/http"
	"text/template"

	"bcintranet/store"
	"bcintranet/templates"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

func GetValidSession(req *http.Request) (*sessions.Session, error) {
	// Returns a valid authenticated user session
	sessStore := sessions.NewCookieStore([]byte("gplus"))
	return sessions.Store.Get(sessStore, req, "gplus_gothic_session")
}

func CustomTemplateExecute(res http.ResponseWriter, req *http.Request, templateName string, data map[string]interface{}) {
	// Append common templates and data structs and execute template
	t, _ := template.ParseFiles(templates.BASE, templates.NOTIFICATIONS, templates.TICKER, templateName)
	if len(data) == 0 {
		data := make(map[string]interface{})
		data["user"], _ = store.GetUser(context.Get(req, "userid").(string))
		t.Execute(res, data)
	} else {
		data["user"], _ = store.GetUser(context.Get(req, "userid").(string))
		t.Execute(res, data)
	}
}
