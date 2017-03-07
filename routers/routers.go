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
	// media route
	common.PathPrefix(urls.MEDIA_PATH).Handler(
		http.StripPrefix(urls.MEDIA_PATH, http.FileServer(http.Dir("media"))))
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
	bc.Get(urls.PAYSLIP_HISTORY_PATH, controllers.PayslipHistoryController)
	bc.Get(urls.PAYSLIP_PATH, controllers.PayslipController)
	bc.Post(urls.PAYSLIP_PATH, controllers.PayslipController)
	bc.Get(urls.USERS_PATH, controllers.UsersController)
	bc.Get(urls.DAILY_LOGS_VIEW_PATH, controllers.DailyLogsViewController)
	bc.Get(urls.DAILY_LOGS_ADD_PATH, controllers.DailyLogsAddController)
	bc.Post(urls.DAILY_LOGS_ADD_PATH, controllers.DailyLogsAddController)
	bc.Get(urls.DAILY_LOGS_PATH, controllers.DailyLogsController)
	bc.Get(urls.METRICS_PATH, controllers.MetricsController)
	bc.Get(urls.ANNOUNCEMENTS_PATH, controllers.AnnouncementsController)
	bc.Post(urls.COMMENT_ADD_PATH, controllers.WallPostCommentAddController)
	bc.Get(urls.COMMENT_DELETE_PATH, controllers.WallPostCommentDeleteController)
	bc.Get(urls.WALL_POST_PATH, controllers.WallPostController)
	bc.Get(urls.POST_ADD_PATH, controllers.WallPostAddController)
	bc.Post(urls.POST_ADD_PATH, controllers.WallPostAddController)
	bc.Get(urls.POST_DELETE_PATH, controllers.WallPostDeleteController)
	bc.Post(urls.ABOUT_FEEDBACK_PATH, controllers.AboutFeedbackController)
	bc.Get(urls.ABOUT_PATH, controllers.AboutController)
	bc.Post(urls.PROJECT_FEEDBACK_PATH, controllers.ProjectFeedbackController)
	bc.Get(urls.PROJECT_VIEW_PATH, controllers.ProjectViewController)
	bc.Get(urls.PROJECT_EDIT_PATH, controllers.ProjectEditController)
	bc.Post(urls.PROJECT_EDIT_PATH, controllers.ProjectEditController)
	bc.Get(urls.PROJECT_ADD_PATH, controllers.ProjectAddController)
	bc.Post(urls.PROJECT_ADD_PATH, controllers.ProjectAddController)
	bc.Get(urls.PROJECTS_PATH, controllers.ProjectsController)
	bc.NotFoundHandler = http.HandlerFunc(controllers.NotFound)
	// admin only routes
	admin := pat.New()
	admin.Get(urls.APPROVALS_PAYSLIPS_PATH, controllers.AdminApprovalPayslipController)
	admin.Post(urls.APPROVALS_PAYSLIPS_PATH, controllers.AdminApprovalPayslipController)
	admin.Get(urls.APPROVALS_PATH, controllers.AdminApprovalController)
	admin.Get(urls.ADMIN_METRICS_OPS_ADD_PATH, controllers.AdminMetricsOpsAddController)
	admin.Post(urls.ADMIN_METRICS_OPS_ADD_PATH, controllers.AdminMetricsOpsAddController)
	admin.Post(urls.CANDIDATE_COMMENT_ADD_PATH, controllers.CandidateCommentAddController)
	admin.Get(urls.CANDIDATE_ADD_PATH, controllers.CandidateAddController)
	admin.Post(urls.CANDIDATE_ADD_PATH, controllers.CandidateAddController)
	admin.Get(urls.CANDIDATE_VIEW_PATH, controllers.CandidateViewController)
	admin.Get(urls.CANDIDATE_EDIT_PATH, controllers.CandidateEditController)
	admin.Post(urls.CANDIDATE_EDIT_PATH, controllers.CandidateEditController)
	admin.Get(urls.RECRUITMENT_PATH, controllers.RecruitmentController)
	admin.NotFoundHandler = http.HandlerFunc(controllers.NotFound)
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
