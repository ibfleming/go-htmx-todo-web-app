package handlers

import (
	"net/http"
	"time"
)

type PostLogoutHandler struct {
	sessionCookie string
}

type PostLogoutHandlerParams struct {
	SessionCookie string
}

func NewPostLogoutHandler(params PostLogoutHandlerParams) *PostLogoutHandler {
	return &PostLogoutHandler{
		sessionCookie: params.SessionCookie,
	}
}

func (h *PostLogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    h.sessionCookie,
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
