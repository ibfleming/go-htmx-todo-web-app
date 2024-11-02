package handlers

import (
	"log"
	"net/http"
)

type PostLoginHandler struct{}

func NewPostLoginHandler() *PostLoginHandler {
	return &PostLoginHandler{}
}

func (h *PostLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	log.Printf("%s %s", email, password)
}
