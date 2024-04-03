package server

import (
	"fmt"
	"net/http"
	"receptionStudio/dbHelper"
	"receptionStudio/handlers"
	"receptionStudio/webSocket"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func StartServer() {
	err := dbHelper.Seed()
	if err != nil {
		fmt.Println("Error :::: ", err)
		return
	}
	http.HandleFunc("/user", handlers.UserHandler)
	http.Handle("/CreateUser", enableCORS(http.HandlerFunc(handlers.CreateUserHandlerStep1)))
	http.Handle("/AuthUser", enableCORS(http.HandlerFunc(handlers.CreateUserHandlerStep2)))
	http.HandleFunc("/Oauth", handlers.OAuthGoogleLogin)
	http.HandleFunc("/OauthCallback", handlers.OAuthCallBack)
	http.Handle("/login", enableCORS(http.HandlerFunc(handlers.LoginWithMailHash)))
	http.HandleFunc("/projectCreate", handlers.CreateProject)
	http.HandleFunc("/projectUpdate", handlers.UpdateProject)
	http.HandleFunc("/projectDelete", handlers.DeleteProject)
	http.HandleFunc("/projectShare", handlers.ShareProject)
	http.HandleFunc("/projectRead", handlers.ReadProject)
	http.Handle("/info", enableCORS(http.HandlerFunc(handlers.ProjectInfo)))
	http.Handle("/websocket", enableCORS(http.HandlerFunc(webSocket.WsEndpoint)))
	fmt.Println("Start Server at localhost:8080")
	http.ListenAndServe(":8080", nil)
}
