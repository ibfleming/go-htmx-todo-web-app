package handlers

import (
	"net/http"
	"zion/internal/storage"
	"zion/internal/utils"
	"zion/templates"
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

	// Check for errors
	if err != nil {
		errMsg := utils.SQLErrorMessage(utils.ExtractSQLStateErrorCode(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		templates.RegisterError(errMsg).Render(r.Context(), w)
		return
	}

	// Render register success template
	c := templates.RegisterSuccess()
	err = c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "❌ Error rendering template", http.StatusInternalServerError)
	}
}
