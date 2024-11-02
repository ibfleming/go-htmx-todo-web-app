package handlers

import (
	"zion/templates"
	"net/http"
)

type GetLoginHandler struct{}

func NewGetLoginHandler() *GetLoginHandler {
	return &GetLoginHandler{}
}

func (h *GetLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.Login()
	err := c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "❌ Error rendering template", http.StatusInternalServerError)
		return
	}
}
