package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"receptionStudio/auth"
	"receptionStudio/dbHelper"
	"strconv"
	"time"
)

type RequestBodyCreateUserStep1 struct {
	Email    string
	Password string
}

func CreateUserHandlerStep1(w http.ResponseWriter, r *http.Request) {
	/*
		Generate jwt & VerifyCode
		Insert Email, Hash, VerifyCode into VerifyDB
		Return jwt
		Send Email with VerifyCode
	*/
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var requestBody RequestBodyCreateUserStep1
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	requestMail := requestBody.Email
	requestHash := requestBody.Password
	if requestHash == "" || requestMail == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hash := sha256.New()
	hash.Write([]byte(requestHash))
	hashed := hash.Sum(nil)
	requestHash = hex.EncodeToString(hashed)
	exists, err := dbHelper.CheckExistUserTable(requestMail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "SOMETHING WENT WRONG")
		return
	}
	if exists {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "EMAIL ALLREADY USED")
	}
	accessToken := auth.CreateUser(requestMail)
	fmt.Fprintf(w, accessToken)
	rand.Seed(time.Now().UnixNano())
	verifyCode := strconv.Itoa(rand.Intn(9999999999))
	auth.SendConfirmEmail(requestMail, verifyCode)
	err = dbHelper.InsertIntoVerifyDB(requestMail, requestHash, verifyCode)
	if err != nil {
		return
	}
}

type RequestBodyCreateUserStep2 struct {
	TokenString string
	VerifyCode  string
}

func CreateUserHandlerStep2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var requestBody RequestBodyCreateUserStep2
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	requestToken := requestBody.TokenString
	requestCode := requestBody.VerifyCode
	mail, err := auth.CheckJwt(requestToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error createHandler.go:90", err)
		fmt.Fprintf(w, "SOMETHING WENT WRONG")
		return
	}
	verify, err := dbHelper.CheckVerifyCode(mail, requestCode)
	if errors.Is(err, fmt.Errorf("JwtIsExpired")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error createHandler.go:101", err)
		fmt.Fprintf(w, "SOMETHING WENT WRONG")
		return
	}
	if verify {
		dbHelper.MoveFromVerifyDBToUserDB(mail)
		returnJwt, err := auth.GenerateJwtAuthGeneral(mail)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Error createHandler.go:110", err)
			fmt.Fprintf(w, "SOMETHING WENT WRONG")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, returnJwt)
	}
}
