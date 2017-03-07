package main

import (
	// system local third-party

	"os"

	"bcintranet/helpers"
	"bcintranet/models"
	"bcintranet/routers"
	"bcintranet/store"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/gplus"
	"github.com/urfave/negroni"
)

func init() {
	// goth package cookie store initialization
	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("bc_app_key")))
}

func main() {
	// use goth provider for authentication
	goth.UseProviders(
		gplus.New(
			os.Getenv("bc_intranet_client_id"),
			os.Getenv("bc_intranet_client_secret"),
			os.Getenv("bc_host")+"/auth/gplus/callback",
		),
	)
	//Insert about data
	_, err := store.GetAbout()
	if err != nil {
		var about models.About
		about = helpers.InsertAboutData()
		store.SaveAbout(&about)
	}
	// get pat router from routers package
	p := routers.GetRouter()
	// use negroni handler
	n := negroni.Classic()
	n.UseHandler(p)
	// run on 3001 and using gin(repl) on 3000
	var port string
	if os.Getenv("bc_env") == "development" {
		port = "3001"
	} else {
		port = os.Getenv("PORT")
	}
	n.Run(":" + port)
}
