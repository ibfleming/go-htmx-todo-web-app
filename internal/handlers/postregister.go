package handlers

import (
	"net/http"
	"zion/internal/storage"
	"zion/internal/utils"
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
		errMsg := utils.SQLErrorMessage(utils.ExtractSQLStateErrorCode(err.Error()))

		// Render register error template
		err = templates.RegisterError(errMsg).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "❌ Error rendering template", http.StatusInternalServerError)
		}
		return
	}

	// Render register success template
	err = templates.RegisterSuccess().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "❌ Error rendering template", http.StatusInternalServerError)
	}
}
