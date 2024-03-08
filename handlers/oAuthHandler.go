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
		fmt.Fprintf(w, "Invalid oauth state")
		return
	}
	code := r.FormValue("code")
	token, err := auth.CodeExchange(code)
	if err != nil {
		fmt.Fprintln(w, "Code exchange failed with error")
		return
	}
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		fmt.Fprintf(w, "Failed to get user info")
		return
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		fmt.Fprintf(w, "Failed to decode JSON")
		return
	}
	fmt.Println(userInfo)
	email, ok := userInfo["email"].(string)
	if !ok {
		fmt.Fprintf(w, "Email not found in response")
		return
	}
	jwtToken, err := auth.GenerateJwtAuthGeneral(email)
	if err != nil {
		fmt.Fprintf(w, "Can't generate JWT")
		return
	}
	exists, err := dbHelper.CheckExistUserTable(email)
	if err != nil {
		fmt.Fprintf(w, "Database Error")
		return
	}
	if !exists {
		dbHelper.InsertIntoUserTable(email, "oauth")
	}
	fmt.Fprintf(w, jwtToken)
}
