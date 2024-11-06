package handlers

import (
	"net/http"
	"zion/internal/middleware/auth"
	"zion/templates"
)

type GetLoginHandler struct{}

func NewGetLoginHandler() *GetLoginHandler {
	return &GetLoginHandler{}
}

func (h *GetLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if auth.IsLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		err := templates.Login().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
			return
		}
	}
}
