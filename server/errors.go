package server

import (
	"log"
	"net/http"
)

func Errors(w http.ResponseWriter, code int, message string) {
	log.Printf("ERROR: code - %d, text: %s", code, message)

	d := struct {
		ErrorCode int
		ErrorText string
	}{
		ErrorCode: code,
		ErrorText: message,
	}

	if err := tpl.ExecuteTemplate(w, "error.html", d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
}
