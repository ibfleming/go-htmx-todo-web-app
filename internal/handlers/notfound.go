package handlers

import (
	"zion/templates"
	"net/http"
)

type GetNotFound struct{}

func NewGetNotFound() *GetNotFound {
	return &GetNotFound{}
}

func (h *GetNotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	c := templates.NotFound()
	err := c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "❌ Error rendering template", http.StatusInternalServerError)
		return
	}
}
