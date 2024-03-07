package server

import (
	"net/http"
	"receptionStudio/handlers"
)

func StartServer() {
	http.HandleFunc("/user", handlers.UserHandler)
	http.HandleFunc("/CreateUser", handlers.CreateUserHandlerStep1)
	http.HandleFunc("/AuthUser", handlers.CreateUserHandlerStep2)
	http.ListenAndServe(":8080", nil)
}
