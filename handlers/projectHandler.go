package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"receptionStudio/auth"
	"receptionStudio/dbHelper"
)

type requestBodyCreate struct {
	ProjectName  string
	RequestToken string
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var requestBody requestBodyCreate
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	requestToken := requestBody.RequestToken
	requestName := requestBody.ProjectName
	if requestName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mail, err := auth.CheckJwt(requestToken)
	if err == fmt.Errorf("JwtIsExpired") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = dbHelper.InsertNewProject(mail, requestName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

type requestBodyUpdate struct {
	ProjectNumber string
	ProjectName   string
	ProjectData   string
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var requestBody requestBodyUpdate
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	requestNum := requestBody.ProjectNumber
	requestName := requestBody.ProjectName
	requestData := requestBody.ProjectData
	err = dbHelper.UpdateOldProject(requestNum, requestName, requestData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
