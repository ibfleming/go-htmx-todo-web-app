package auth

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	zerr "zion/internal/errors"
	"zion/internal/storage"
	"zion/internal/storage/schema"
)

type AuthMiddleware struct {
	sessions          storage.SessionStorageInterface
	sessionCookieName string
}

type AuthMiddlewareParams struct {
	Sessions          storage.SessionStorageInterface
	SessionCookieName string
}

func NewAuthMiddleware(params AuthMiddlewareParams) *AuthMiddleware {
	return &AuthMiddleware{
		sessions:          params.Sessions,
		sessionCookieName: params.SessionCookieName,
	}
}

type UserContextKey string

var UserKey UserContextKey = "user"

func (m *AuthMiddleware) AddUserToContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(m.sessionCookieName)
		if err != nil {
			log.Print(zerr.ErrInvalidCookie.Error())
			h.ServeHTTP(w, r)
			return
		}
		decodedValue, err := base64.StdEncoding.DecodeString(sessionCookie.Value)
		if err != nil {
			log.Print(zerr.ErrFailedToDecodeString.Error())
			h.ServeHTTP(w, r)
			return
		}
		splitValue := strings.Split(string(decodedValue), ":")
		if len(splitValue) != 2 {
			log.Print(zerr.ErrInvalidCookieFormat.Error())
			h.ServeHTTP(w, r)
			return
		}
		sessionID := splitValue[0]
		userID := splitValue[1]
		user, err := m.sessions.GetUserFromSession(sessionID, userID)
		if err != nil {
			log.Print(zerr.ErrFailedToGetUserFromSession.Error())
			h.ServeHTTP(w, r)
			return
		}
		log.Printf("[AUTH] Authenticated - UserID: %d, Email: %s", user.ID, user.Email)
		ctx := context.WithValue(r.Context(), UserKey, user)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(ctx context.Context) *schema.User {
	user := ctx.Value(UserKey)
	if user == nil {
		return nil
	}
	return user.(*schema.User)
}

func IsLoggedIn(r *http.Request) bool {
	return GetUser(r.Context()) != nil
}
