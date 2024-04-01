package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"receptionStudio/auth"
	"receptionStudio/dbHelper"
)

type requestLoginBody struct {
	Email    string
	Password string
}

func LoginWithMailHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var requestBody requestLoginBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "SOMETHING WENT WRONG")
		return
	}
	requestMail := requestBody.Email
	requestHash := requestBody.Password
	hash := sha256.New()
	hash.Write([]byte(requestHash))
	hashed := hash.Sum(nil)
	requestHash = hex.EncodeToString(hashed)
	fmt.Println("RequestMail :", requestMail)
	fmt.Println("RequestHash :", requestHash)
	check, err := dbHelper.CheckMailWithHash(requestMail, requestHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "SOMETHING WENT WRONG")
		return
	}
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	accessToken, err := auth.GenerateJwtAuthGeneral(requestMail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "SOMETHING WENT WRONG")
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, accessToken)
}
