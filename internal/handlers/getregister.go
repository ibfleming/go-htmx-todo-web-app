package handlers

import (
	"net/http"
	"zion/internal/middleware/auth"
	"zion/templates"
)

type GetRegisterHandler struct{}

func NewGetRegisterHandler() *GetRegisterHandler {
	return &GetRegisterHandler{}
}

func (h *GetRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if auth.IsLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		err := templates.Register().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
			return
		}
	}
}
