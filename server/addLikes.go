package server

import (
	"fmt"
	"net/http"

	"forum/models"
	"forum/service"
)

// memberLikeForPost
func (h *Handler) memberLikeForPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/post/like" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		// TODO add context err message
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	currentUserId := user.Id
	var like models.LikePost
	like.AuthorID = user.Id
	like.PostID = r.FormValue("postId")
	like.Type = r.FormValue("type") == "true"

	_, err := service.AddLikePost(h.repos, like, currentUserId)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post-page?id=%s", like.PostID), http.StatusFound)
}

// memberLikeForComment
func (h *Handler) memberLikeForComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/comment/like" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		// TODO add context err message
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	currentUserId := user.Id
	postIdPage := r.FormValue("postId")

	var like models.LikeComment
	like.AuthorID = user.Id
	like.CommentID = r.FormValue("commentId")
	like.Type = r.FormValue("type") == "true"

	_, err := service.AddLikeComment(h.repos, like, currentUserId)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post-page?id=%s", postIdPage), http.StatusFound)
}
