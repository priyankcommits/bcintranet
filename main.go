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
	// goth package cookie store initialization
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
	// get pat router from routers package
	p := routers.GetRouter()
	// use negroni handler
	n := negroni.Classic()
	n.UseHandler(p)
	// run on 3001 and using gin(repl) on 3000
	port := os.Getenv("PORT")
	n.Run(":" + port)
}
