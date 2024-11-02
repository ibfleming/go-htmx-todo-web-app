package handlers

import (
	"zion/internal/storage"
	"zion/templates"
	"log"
	"net/http"
)

type PostRegisterHandler struct {
	users storage.UserStorage
}

type PostRegisterHandlerParameters struct {
	Users storage.UserStorage
}

func NewPostRegisterHandler(params PostRegisterHandlerParameters) *PostRegisterHandler {
	return &PostRegisterHandler{
		users: params.Users,
	}
}

func (h *PostRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := h.users.CreateUser(email, password)

	// Error 2
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		templates.RegisterError(err.Error()).Render(r.Context(), w)
		return
	}

	// Render register success template
	c := templates.RegisterSuccess()
	err = c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "❌ Error rendering template", http.StatusInternalServerError)
	}

	log.Printf("📝 User registered succssfully!\n\tEmail: %s\n\tPassword: %s", email, password)
}
