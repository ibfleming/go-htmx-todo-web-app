package server

import (
	"net/http"
	"zion/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func CreateRouter(s *Server) http.Handler {
	// Create a new chi router
	mux := chi.NewRouter()

	// Initialize file server
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	mux.Handle("/static/*", fs)
	mux.Handle("/favicon.ico", fs)

	mux.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
			middleware.Recoverer,
			middleware.NoCache,
			middleware.Heartbeat("/ping"),
		)

		r.NotFound(handlers.NewGetNotFound().ServeHTTP)

		r.Get("/", handlers.NewGetNotFound().ServeHTTP)

		r.Get("/login", handlers.NewGetLoginHandler().ServeHTTP)
		r.Post("/login", handlers.NewPostLoginHandler(handlers.PostLoginHandlerParameters{
			Users:         *s.users,
			Sessions:      *s.sessions,
			PasswordHash:  *s.hash,
			SessionCookie: s.sessionCookie,
		}).ServeHTTP)

		r.Get("/register", handlers.NewGetRegisterHandler().ServeHTTP)
		r.Post("/register", handlers.NewPostRegisterHandler(handlers.PostRegisterHandlerParameters{
			Users: *s.users,
		}).ServeHTTP)
	})

	return mux
}
