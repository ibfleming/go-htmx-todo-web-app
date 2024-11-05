package handlers

import (
	"net/http"
	"zion/internal/middleware/auth"
	"zion/templates"

	"github.com/a-h/templ"
)

type Index struct{}

func NewIndex() *Index {
	return &Index{}
}

func (h *Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	var c templ.Component
	if user != nil {
		c = templates.Index(user)
	} else {
		w.WriteHeader(http.StatusNotFound)
		c = templates.NotFound()
	}

	err := c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "❌ Error rendering template", http.StatusInternalServerError)
		return
	}
}
