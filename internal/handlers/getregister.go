package handlers

import (
	"zion/templates"
	"net/http"
)

type GetRegisterHandler struct{}

func NewGetRegisterHandler() *GetRegisterHandler {
	return &GetRegisterHandler{}
}

func (h *GetRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.Register()
	err := c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "❌ Error rendering template", http.StatusInternalServerError)
		return
	}
}
