package server

import (
	"net/http"

	"forum/models"
	"forum/service"
)

// getAllPosts handler - GET only
func (h *Handler) getAllPosts(w http.ResponseWriter, r *http.Request) {
	//---negative cases---
	if r.URL.Path != "/posts" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	responseBody := struct {
		Posts  []models.Post `json:"posts"`
		IsAuth bool          `json:"autorized"`
		UserId int           `json:"userId"`
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

	if tokenIsValid {
		responseBody.IsAuth = true
		responseBody.UserId = user.Id
		allPosts, err := service.GetAllPosts(h.repos, user.Id)
		if err != nil {
			Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		responseBody.Posts = allPosts
	} else {
		responseBody.IsAuth = false
		responseBody.UserId = 0
		allPosts, err := service.GetAllPosts(h.repos, 0)
		if err != nil {
			Errors(w, http.StatusInternalServerError, err.Error())
			return
		}

		responseBody.Posts = allPosts
	}

	if err := tpl.ExecuteTemplate(w, "allPosts.html", responseBody); err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
}
