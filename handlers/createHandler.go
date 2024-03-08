package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"receptionStudio/auth"
	"receptionStudio/dbHelper"
	"strconv"
	"time"
)

type RequestBodyCreateUserStep1 struct {
	Mail string
	Hash string
}

func CreateUserHandlerStep1(w http.ResponseWriter, r *http.Request) {
	/*
		Generate jwt & VerifyCode
		Insert Email, Hash, VerifyCode into VerifyDB
		Return jwt
		Send Email with VerifyCode
	*/
	if r.Method != http.MethodPost {
		return
	}
	var requestBody RequestBodyCreateUserStep1
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		return
	}
	requestMail := requestBody.Mail
	requestHash := requestBody.Hash
	exists, err := dbHelper.CheckExistUserTable(requestMail)
	if err != nil {
		fmt.Fprintf(w, "Database Error")
		return
	}
	if exists {
		fmt.Fprintf(w, "Mail Address is already used.")
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
		return
	}
	var requestBody RequestBodyCreateUserStep2
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		return
	}
	requestToken := requestBody.TokenString
	requestCode := requestBody.VerifyCode
	mail, err := auth.CheckJwt(requestToken)
	if err != nil {
		fmt.Println("ERROR :", err)
		//fmt.Println("REQUEST TOKEN :", requestToken)
		//fmt.Println("REQUEST CODE :", requestCode)
		fmt.Fprintf(w, "Can't Auth")
		return
	}
	verify, err := dbHelper.CheckVerifyCode(mail, requestCode)
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}
	if verify {
		returnJwt, err := auth.GenerateJwtAuthGeneral(mail)
		if err != nil {
			fmt.Fprintf(w, "Can't generate jwt")
			return
		}
		fmt.Fprintf(w, returnJwt)
	}
}
