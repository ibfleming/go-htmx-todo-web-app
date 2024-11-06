package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
	"zion/internal/errors"
	"zion/internal/hash"
	"zion/internal/storage"
	"zion/internal/storage/db"
	"zion/templates"
)

type PostLoginHandler struct {
	users         storage.UserStorageInterface
	sessions      storage.SessionStorageInterface
	passwordHash  *hash.PasswordHash
	sessionCookie string
}

type PostLoginHandlerParameters struct {
	Users         storage.UserStorageInterface
	Sessions      storage.SessionStorageInterface
	PasswordHash  *hash.PasswordHash
	SessionCookie string
}

func NewPostLoginHandler(params PostLoginHandlerParameters) *PostLoginHandler {
	return &PostLoginHandler{
		users:         params.Users,
		sessions:      params.Sessions,
		passwordHash:  params.PasswordHash,
		sessionCookie: params.SessionCookie,
	}
}

func (h *PostLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	// Get the user from the database
	user, err := h.users.GetUser(email)
	// Error handling
	if err != nil {
		// Error header
		w.WriteHeader(http.StatusUnauthorized)
		// Render register error template
		err = templates.LoginError(errors.ErrInvalidEmailOrPassword.Error()).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
		}
		return
	}
	// Check if the password is valid matches the hash
	isPasswordValid, err := h.passwordHash.ComparePasswordAndHash(password, user.Password)
	// Error handling
	if err != nil || !isPasswordValid {
		// Error header
		w.WriteHeader(http.StatusUnauthorized)
		// Render register error template
		err = templates.LoginError(errors.ErrPasswordIncorrect.Error()).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
		}
		return
	}
	// Create the session since the user is valid
	session, err := h.sessions.CreateSession(&db.Session{
		UserID: user.ID,
	})
	// Error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Cookie Value (sessionID:userID)
	cookieValue := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%d", session.SessionID, user.ID)))
	// Create session cookie
	cookie := http.Cookie{
		Name:     h.sessionCookie,
		Value:    cookieValue,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	// Redirect
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
