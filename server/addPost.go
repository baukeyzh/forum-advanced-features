package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"forum/models"
	"forum/service"
)

func (h *Handler) memberPostCreate(w http.ResponseWriter, r *http.Request) {
	var activity models.Activity
	if r.URL.Path != "/v1/post/create" {
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
	imageUrl := "./imgs/"
	timestamp := time.Now().String()
	var post models.Post
	id := r.FormValue("post-id")
	idInt, _ := strconv.Atoi(id)
	post.Id = idInt
	if post.Id == 0 {
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		file, handler, err := r.FormFile("postImage")
		if err != nil {
			Errors(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(handler.Filename))

		if ext != ".jpeg" && ext != ".jpg" && ext != ".gif" && ext != ".svg" && ext != ".png" {
			Errors(w, http.StatusBadRequest, ext+"File type not supported. Only .jpg, .jpeg, .gif, .svg, .png files are allowed.")
			return
		}
		if handler.Size > 5*1024*1024 {
			Errors(w, http.StatusBadRequest, "File size exceeds the limit of 10 MB")
			return
		}
		timestamp := time.Now().String()
		imageName := timestamp + ext
		// Create a new file on the server to save the image
		f, err := os.OpenFile(imageUrl+imageName, os.O_WRONLY|os.O_CREATE, 0o666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		// Copy the image to the new file
		io.Copy(f, file)
		post.ImageName = imageName
	} else {
		err := r.ParseMultipartForm(1 << 20) // 10 MB
		if err != nil {
			Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		file, handler, _ := r.FormFile("postImage")
		if file != nil {
			ext := filepath.Ext(handler.Filename)
			if ext != ".jpeg" && ext != ".jpg" && ext != ".gif" && ext != ".svg" && ext != ".png" {
				Errors(w, http.StatusBadRequest, "File type not supported. Only .jpg, .jpeg, .gif, .svg, .png files are allowed.")
				return
			}
			if handler.Size > 5*1024*1024 {
				Errors(w, http.StatusBadRequest, "File size exceeds the limit of 10 MB")
				return
			}
			imageName := timestamp + ext
			defer file.Close()
			f, err := os.OpenFile(imageUrl+imageName, os.O_WRONLY|os.O_CREATE, 0o666)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			io.Copy(f, file)
			post.ImageName = imageName
		} else {
			postDb, _ := service.GetPostById(h.repos, post.Id)
			post.ImageName = postDb.ImageName
		}
	}
	post.Title = r.FormValue("postTitle")
	post.Content = r.FormValue("postContent")
	post.AuthorName = user.UserName
	post.AuthorID = user.Id
	var cats []int
	if r.FormValue("1") == "on" {
		cats = append(cats, 1)
	}
	if r.FormValue("2") == "on" {
		cats = append(cats, 2)
	}
	if r.FormValue("3") == "on" {
		cats = append(cats, 3)
	}
	if r.FormValue("4") == "on" {
		cats = append(cats, 4)
	}

	if len(cats) == 0 {
		cats = []int{1, 2, 3, 4}
	}
	if post.Id == 0 {
		_, err := service.AddPost(h.repos, post, cats, user.Id)
		if err != nil {
			Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		activity.Action = models.Create
	} else {
		editedPost, err := service.SetPost(h.repos, post, cats, currentUserId)
		if err != nil {
			Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		post = editedPost
		activity.Action = models.Edit
	}
	// TODO add context of error or created post
	http.Redirect(w, r, "/posts", http.StatusFound)
}

func (h *Handler) memberPostEdit(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/post/edit" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	id := r.FormValue("post-id")
	idInt, _ := strconv.Atoi(id)
	post, err := service.GetPostById(h.repos, idInt)
	if err != nil {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}
	_, ok := r.Context().Value("user").(models.User)
	if !ok {
		// TODO add context err message
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if err := tpl.ExecuteTemplate(w, "edit_post.html", post); err != nil {
		Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
	return
}

func (h *Handler) memberPostDelete(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/post/delete" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	id := r.FormValue("post-id")
	idInt, _ := strconv.Atoi(id)
	_, ok := r.Context().Value("user").(models.User)
	if !ok {
		// TODO add context err message
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	err := service.DeletePostById(h.repos, idInt)
	if err != nil {
		Errors(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/posts"), http.StatusFound)
}
