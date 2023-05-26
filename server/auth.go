package server

import (
	"encoding/json"
	"net/http"
	"time"

	"forum/models"
	"forum/service"
)

const CookieName = "token"

// gestRegistration handler -GET/POST
func (h *Handler) gestRegistration(w http.ResponseWriter, r *http.Request) {
	//---negative cases---
	if r.URL.Path != "/registration" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	//---positive cases---
	// GET /registration
	if r.Method == http.MethodGet {
		errorMessage := r.FormValue("error")
		var d struct {
			ErrorMessage   string
			SuccessMessage string
		}
		if errorMessage != "" {
			d = struct {
				ErrorMessage   string
				SuccessMessage string
			}{
				ErrorMessage:   errorMessage,
				SuccessMessage: "hello Alem",
			}
		}

		if err := tpl.ExecuteTemplate(w, "registration.html", d); err != nil {
			Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		return
	}

	// POST /registration
	if r.Method == http.MethodPost {
		userName := r.FormValue("registerUsername")
		email := r.FormValue("registerEmail")
		password := r.FormValue("registerPassword")
		confirmPassword := r.FormValue("registerConfirmPassword")
		if password != confirmPassword {
			http.Redirect(w, r, "/registration?error=Password is not match", http.StatusBadRequest)
			return
		}

		code, err := service.Registration(h.repos, userName, email, password)
		if err != nil || code != http.StatusCreated {
			http.Redirect(w, r, "/registration?error=User name and email must be unique", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
}

// gestLogin handler -POST only
func (h *Handler) gestLogin(w http.ResponseWriter, r *http.Request) {
	//---negative cases---
	if r.URL.Path != "/login" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	//---positive cases---
	// GET /login
	if r.Method == http.MethodGet {
		errorMessage := r.FormValue("error")
		var d struct {
			ErrorMessage string
		}
		if errorMessage != "" {
			d = struct {
				ErrorMessage string
			}{
				ErrorMessage: errorMessage,
			}
		}
		if err := tpl.ExecuteTemplate(w, "login.html", d); err != nil {
			Errors(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// POST /login
	if r.Method == http.MethodPost {
		email := r.FormValue("loginEmail")
		password := r.FormValue("loginPassword")
		user, code, err := service.Authentication(h.repos, email, password)
		res := struct {
			Value    string    `json:"token"`
			Expires  time.Time `json:"expires"`
			ErrorMsg string    `json:"error"`
		}{
			Value:   "",
			Expires: user.ExpireAt,
		}
		if err != nil {
			if code == http.StatusUnauthorized || code == http.StatusBadRequest {
				w.WriteHeader(code)
			} else {
				Errors(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		token, err := service.Authorization(h.repos, user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		res.Value = token
		jsonRes, err := json.Marshal(res)
		if err != nil {
			Errors(w, http.StatusBadRequest, err.Error())
			return
		}
		w.Write([]byte(jsonRes))
	}

	//	http.SetCookie(w, c)

	//	http.Redirect(w, r, "/posts", http.StatusSeeOther)
	return
}

// memberLogout handler -GET only
func (h *Handler) memberLogout(w http.ResponseWriter, r *http.Request) {
	//---negative cases---
	if r.URL.Path != "/logout" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	//---positive cases---
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if err := service.Logout(h.repos, user); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	c := &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, c)

	http.Redirect(w, r, "/", http.StatusFound)
}
