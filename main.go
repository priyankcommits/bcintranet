package main

import (
	"os"

	"bcintranet/routers"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/gplus"
	"github.com/urfave/negroni"
)

func init() {
	gothic.Store = sessions.NewCookieStore([]byte("gplus"))
}

func main() {
	// use goth provider for authentication
	goth.UseProviders(
		gplus.New(
			os.Getenv("bc_intranet_client_id"),
			os.Getenv("bc_intranet_client_secret"),
			"http://localhost:3000/auth/gplus/callback",
		),
	)
	// get pat router
	p := routers.GetRouter()
	// use negroni middleware
	n := negroni.Classic()
	n.UseHandler(p)
	// run on 3001 and using gin on 3000
	n.Run(":3001")
}
