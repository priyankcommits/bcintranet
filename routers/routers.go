package routers

import (
	"net/http"

	"bcintranet/controllers"
	"bcintranet/middlewares"
	"bcintranet/urls"

	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
)

// registers all routes for the application.
func GetRouter() *pat.Router {
	// url paths imported from urls package
	// common routes
	common := pat.New()
	// static route
	common.PathPrefix(urls.STATIC_PATH).Handler(
		http.StripPrefix(urls.STATIC_PATH, http.FileServer(http.Dir("static"))))
	common.Get(urls.AUTH_CALLBACK_PATH, controllers.AuthCallbackController)
	common.Get(urls.AUTH_PATH, controllers.AuthController)
	common.Get(urls.LOGOUT_PATH, controllers.LogoutController)
	common.Get(urls.NOT_FOUND_PATH, controllers.NotFound)
	// Profile routes
	profile := pat.New()
	profile.Get(urls.PROFILE_PATH_VIEW, controllers.ProfileViewController)
	profile.Get(urls.PROFILE_PATH_EDIT, controllers.ProfileEditController)
	profile.Post(urls.PROFILE_PATH_EDIT, controllers.ProfileEditController)
	// bc routes
	bc := pat.New()
	bc.Get(urls.HOME_PATH, controllers.HomeController)
	// applying middlewares
	common.PathPrefix(urls.PROFILE_PATH).Handler(
		negroni.New(
			negroni.HandlerFunc(
				middlewares.GothLoginMiddleware),
			negroni.Wrap(profile),
		),
	)
	common.PathPrefix(urls.BC_PATH).Handler(
		negroni.New(
			negroni.HandlerFunc(
				middlewares.GothLoginMiddleware),
			negroni.HandlerFunc(
				middlewares.SetUserMiddleware),
			negroni.Wrap(bc),
		),
	)
	common.Get(urls.ROOT_PATH, controllers.LoginController)
	return common
}
