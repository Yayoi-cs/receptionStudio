package server

import (
	"net/http"
	"receptionStudio/handlers"
)

func StartServer() {
	http.HandleFunc("/user", handlers.UserHandler)
	http.HandleFunc("/CreateUser", handlers.CreateUserHandlerStep1)
	http.HandleFunc("/AuthUser", handlers.CreateUserHandlerStep2)
	http.HandleFunc("/Oauth", handlers.OAuthGoogleLogin)
	http.HandleFunc("/OauthCallback", handlers.OAuthCallBack)
	http.HandleFunc("/login", handlers.LoginWithMailHash)
	http.HandleFunc("/projectCreate", handlers.CreateProject)
	http.HandleFunc("/projectUpdate", handlers.UpdateProject)
	http.HandleFunc("/projectDelete", handlers.DeleteProject)
	http.ListenAndServe(":8080", nil)
}
