package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"forum/models"
	"forum/service"
)

func (h *Handler) getActivities(w http.ResponseWriter, r *http.Request) {
	//---negative cases---
	if r.URL.Path != "/activity" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	fmt.Println()
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	var user models.User
	tokenIsValid := false
	token, err := r.Cookie(CookieName)
	if err == nil {
		user, err = service.Identification(h.repos, token.Value)
		if err == nil {
			tokenIsValid = true
		}
	}
	if !tokenIsValid {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	activityPage := struct {
		Activities   []models.Activity `json:"activities"`
		IsAuthorized bool              `json:"isAuthorized"`
	}{}
	activities, err := service.GetActivitiesByCurrentUserId(h.repos, user.Id)
	for i := 0; i < len(activities); i++ {
		t := activities[i].CreatedAt.Format("2006-01-02 15:04:05")
		activities[i].CreatedAtStr = t
	}
	if err != nil {
		Errors(w, http.StatusNotFound, err.Error())
		return
	}
	activityPage.Activities = activities
	activityPage.IsAuthorized = tokenIsValid
	if err := tpl.ExecuteTemplate(w, "activities.html", activityPage); err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// getAllPosts handler - GET only
func (h *Handler) getActivitiesCount(w http.ResponseWriter, r *http.Request) {
	//---negative cases---
	if r.URL.Path != "/get-activities-count" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	responseBody := struct {
		ActivitiesCount int `json:"activitiesCount"`
	}{}
	tokenIsValid := false
	var user models.User
	token, err := r.Cookie(CookieName)
	if err == nil {
		user, err = service.Identification(h.repos, token.Value)
		if err == nil {
			tokenIsValid = true
		}
	}
	activitiesCount, err := service.GetActivitiesCountByCurrentUserId(h.repos, user.Id)
	if tokenIsValid {
		responseBody.ActivitiesCount = activitiesCount
	} else {
		responseBody.ActivitiesCount = 0
	}

	// Convert the struct to JSON
	jsonData, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	w.Write(jsonData)
}
