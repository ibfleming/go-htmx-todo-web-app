package handlers

import (
	"net/http"
	"zion/internal/middleware/auth"
	"zion/templates"
)

type Index struct{}

func NewGetIndexHandler() *Index {
	return &Index{}
}

func (h *Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	err := templates.Index(user).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}
