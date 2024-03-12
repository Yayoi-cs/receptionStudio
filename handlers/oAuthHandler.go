package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"receptionStudio/auth"
	"receptionStudio/dbHelper"
)

func OAuthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := auth.OAuthLoginURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func OAuthCallBack(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if !auth.CheckStateString(state) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "SOMETHING WENT WRONG")
		return
	}
	code := r.FormValue("code")
	token, err := auth.CodeExchange(code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "SOMETHING WENT WRONG")
		return
	}
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "SOMETHING WENT WRONG")
		return
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "SOMETHING WENT WRONG")
		return
	}
	email, ok := userInfo["email"].(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "SOMETHING WENT WRONG")
		return
	}
	jwtToken, err := auth.GenerateJwtAuthGeneral(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "SOMETHING WENT WRONG")
		return
	}
	exists, err := dbHelper.CheckExistUserTable(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "SOMETHING WENT WRONG")
		return
	}
	if !exists {
		dbHelper.InsertIntoUserTable(email, "", true)
	}
	fmt.Fprintf(w, jwtToken)
}
