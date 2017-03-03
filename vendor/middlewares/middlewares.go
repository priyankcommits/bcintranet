package middlewares

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

func GothLoginMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	log.Println("Triggered MIddleware")
	var Store sessions.Store
	session, _ := Store.Get(req, "gplus_gothic_session")

	value := session.Values["gplus"]
	if value == nil {
		log.Println("could not find session")
		// return "", errors.New("could not find a matching session for this request")
	}
	next(res, req)
	log.Println("trigger finished")
}
