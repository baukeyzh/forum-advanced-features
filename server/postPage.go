package server

import (
	"net/http"
	"strconv"

	"forum/models"
	"forum/service"
)

// getOnePostAndComments handler - GET only
// Query selectors: id={int}
func (h *Handler) getPostAndComments(w http.ResponseWriter, r *http.Request) {
	//---negative cases---
	if r.URL.Path != "/post-page" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	//---positive cases---
	postId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || postId == 0 {
		Errors(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	post, err := service.GetPostById(h.repos, postId)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}

	comments, err := service.GetAllComments(h.repos, postId)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}

	postAndComments := models.PostAndComments{
		Post_info: post,
		Comments:  comments,
	}

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
		postAndComments.IsAuth = true
		postAndComments.UserId = user.Id
	} else {
		postAndComments.IsAuth = false
		postAndComments.UserId = 0
	}

	if err := tpl.ExecuteTemplate(w, "postPage.html", postAndComments); err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
}
