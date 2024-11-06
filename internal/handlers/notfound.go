package handlers

import (
	"net/http"
	"zion/templates"
)

type GetNotFound struct{}

func NewGetNotFound() *GetNotFound {
	return &GetNotFound{}
}

func (h *GetNotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	err := templates.NotFound().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "❌ Error rendering template", http.StatusInternalServerError)
		return
	}
}
