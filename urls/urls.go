package urls

//Define route paths to be used in routers package
const STATIC_PATH string = "/static/"
const ROOT_PATH string = "/"
const LOGOUT_PATH string = "/logout"
const NOT_FOUND_PATH string = "/404"
const AUTH_PATH string = "/auth/{provider}"
const AUTH_CALLBACK_PATH string = "/auth/{provider}/callback"
const PROFILE_PATH string = "/profile/"
const PROFILE_PATH_VIEW string = PROFILE_PATH + "view"
const PROFILE_PATH_EDIT string = PROFILE_PATH + "edit"
const BC_PATH = "/bc/"
const HOME_PATH string = BC_PATH + "home"
