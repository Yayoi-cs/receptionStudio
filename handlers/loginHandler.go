package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"receptionStudio/auth"
	"receptionStudio/dbHelper"
)

type requestLoginBody struct {
	Email string
	Hash  string
}

func LoginWithMailHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
	requestHash := requestBody.Hash
	check, err := dbHelper.CheckMailWithHash(requestMail, requestHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "SOMETHING WENT WRONG")
		return
	}
	if !check {
		w.WriteHeader(http.StatusBadRequest)
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
