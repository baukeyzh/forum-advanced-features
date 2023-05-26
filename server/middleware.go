package server

import (
	"context"
	"errors"
	"net/http"

	"forum/models"
	"forum/service"
)

func (h *Handler) identification(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(CookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			Errors(w, http.StatusInternalServerError, err.Error())
			return
		}

		user, err := service.Identification(h.repos, token.Value)
		if err != nil {
			if errors.Is(err, models.ErrorUnauthorized) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		ctx := context.WithValue(context.Background(), "user", user)
		r = r.WithContext(ctx)

		next(w, r)
	})
}
