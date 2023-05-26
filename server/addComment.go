package server

import (
	"fmt"
	"net/http"
	"strconv"

	"forum/models"
	"forum/service"
)

func (h *Handler) memberCommentCreate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/comment/create" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if r.Method != http.MethodPost {
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
	id := r.FormValue("comment-id")
	idInt, _ := strconv.Atoi(id)
	var comm models.Comment
	comm.Content = r.FormValue("comment-text")
	comm.PostID = r.FormValue("post-id")
	comm.Id = idInt
	comm.AuthorName = user.UserName
	comm.AuthorID = user.Id
	if idInt == 0 {
		_, err := service.AddComment(h.repos, comm, currentUserId)
		if err != nil {
			Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		_, err := service.EditComment(h.repos, comm, currentUserId)
		if err != nil {
			Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// w.Write([]byte(fmt.Sprintf("%d", id)))
	// w.WriteHeader(http.StatusCreated)

	http.Redirect(w, r, fmt.Sprintf("/post-page?id=%s", comm.PostID), http.StatusFound)
}

func (h *Handler) memberCommentEdit(w http.ResponseWriter, r *http.Request) {
	var comm models.Comment
	if r.URL.Path != "/v1/comment/edit" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user, ok := r.Context().Value("user").(models.User)
	id := r.FormValue("comment_id")
	idInt, _ := strconv.Atoi(id)
	comm.Id = idInt
	if !ok {
		// TODO add context err message
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	comm, err := service.GetCommentById(h.repos, comm.Id)
	if user.Id != comm.AuthorID {
		Errors(w, http.StatusNotFound, "Not Found")
		return
	}
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
	}
	if err := tpl.ExecuteTemplate(w, "edit_comment.html", comm); err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) memberCommentDelete(w http.ResponseWriter, r *http.Request) {
	var comm models.Comment
	if r.URL.Path != "/v1/comment/delete" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user, ok := r.Context().Value("user").(models.User)
	id := r.FormValue("comment_id")
	idInt, _ := strconv.Atoi(id)
	comm.Id = idInt
	comment, err := service.GetCommentById(h.repos, comm.Id)
	if !ok {
		// TODO add context err message
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if user.Id != comment.AuthorID {
		Errors(w, http.StatusNotFound, "Not Found")
		return
	}
	err = service.DeleteCommentById(h.repos, comm.Id)
	fmt.Println(err)
	if err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
	}
	http.Redirect(w, r, fmt.Sprintf("/post-page?id=%s", comment.PostID), http.StatusFound)
}
