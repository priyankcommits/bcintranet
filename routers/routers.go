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
	common := pat.New()
	// static route
	common.PathPrefix(urls.STATIC_PATH).Handler(
		http.StripPrefix(urls.STATIC_PATH, http.FileServer(http.Dir("static"))))
	// common routes
	common.Get(urls.AUTH_CALLBACK_PATH, controllers.AuthCallbackController)
	common.Get(urls.AUTH_PATH, controllers.AuthController)
	common.Get(urls.LOGOUT_PATH, controllers.LogoutController)
	common.Get(urls.NOT_FOUND_PATH, controllers.NotFound)
	// Profile routes
	profile := pat.New()
	profile.Get(urls.PROFILE_EDIT_PATH, controllers.ProfileEditController)
	profile.Post(urls.PROFILE_EDIT_PATH, controllers.ProfileEditController)
	profile.Get(urls.PROFILE_VIEW_PATH, controllers.ProfileViewController)
	profile.Post(urls.PROFILE_VIEW_PATH, controllers.ProfileViewController)
	profile.Get(urls.HOME_PATH, controllers.HomeController)
	profile.NotFoundHandler = http.HandlerFunc(controllers.NotFound)
	// bc routes
	bc := pat.New()
	bc.Get(urls.PAY_SLIP_PATH, controllers.PaySlipController)
	bc.Get(urls.USERS_PATH, controllers.UsersController)
	bc.Get(urls.METRICS_PATH, controllers.MetricsController)
	bc.NotFoundHandler = http.HandlerFunc(controllers.NotFound)
	// admin only routes
	admin := pat.New()
	admin.Get(urls.ADMIN_METRICS_OPS_ADD_PATH, controllers.AdminMetricsOpsAddController)
	admin.Post(urls.ADMIN_METRICS_OPS_ADD_PATH, controllers.AdminMetricsOpsAddController)
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
	common.PathPrefix(urls.ADMIN_PATH).Handler(
		negroni.New(
			negroni.HandlerFunc(
				middlewares.GothLoginMiddleware),
			negroni.HandlerFunc(
				middlewares.SetUserMiddleware),
			negroni.HandlerFunc(
				middlewares.AdminCheckMiddleware),
			negroni.Wrap(admin),
		),
	)
	common.NotFoundHandler = http.HandlerFunc(controllers.NotFound)
	common.Get(urls.ROOT_PATH, controllers.LoginController)
	return common
}
