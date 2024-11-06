package handlers

import (
	"net/http"
	"zion/internal/errors"
	"zion/internal/storage"
	"zion/templates"
)

type PostRegisterHandler struct {
	users storage.UserStorageInterface
}

type PostRegisterHandlerParameters struct {
	Users storage.UserStorageInterface
}

func NewPostRegisterHandler(params PostRegisterHandlerParameters) *PostRegisterHandler {
	return &PostRegisterHandler{
		users: params.Users,
	}
}

func (h *PostRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	// Create user
	err := h.users.CreateUser(email, password)
	// Error handling
	if err != nil {
		// Error header
		w.WriteHeader(http.StatusBadRequest)
		// Extract SQL state error code
		sqlErr := errors.SQLErrorMessage(err)
		// Render register error template
		err = templates.RegisterError(sqlErr).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
		}
		return
	}
	// Render register success template
	err = templates.RegisterSuccess().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
	}
}
