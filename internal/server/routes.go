package server

import (
	"net/http"
	"zion/internal/handlers"
	"zion/internal/middleware/auth"

	"github.com/go-chi/chi/v5/middleware"
)

/* func (s *ZionServer) EstablishRoutes() {

	authentication := auth.NewAuthMiddleware(auth.AuthMiddlewareParams{
		Sessions:          s.sessions,
		SessionCookieName: s.sessionCookie,
	})

	s.router.Group(func(r chi.Router) {
		// Middleware
		r.Use(
			authentication.AddUserToContext,
			middleware.Logger,
			middleware.Recoverer,
			middleware.NoCache,
			middleware.Heartbeat("/ping"),
		)

		// 404 Handler
		r.NotFound(handlers.NewGetNotFound().ServeHTTP)

		// Home Handler
		r.Get("/", handlers.NewIndex().ServeHTTP)

		// Login Handler
		r.Get("/login", handlers.NewGetLoginHandler().ServeHTTP)
		r.Post("/login", handlers.NewPostLoginHandler(handlers.PostLoginHandlerParameters{
			Users:         s.users,
			Sessions:      s.sessions,
			PasswordHash:  s.hash,
			SessionCookie: s.sessionCookie,
		}).ServeHTTP)

		// Register Handler
		r.Get("/register", handlers.NewGetRegisterHandler().ServeHTTP)
		r.Post("/register", handlers.NewPostRegisterHandler(handlers.PostRegisterHandlerParameters{
			Users: s.users,
		}).ServeHTTP)

		// Logout Handler
		r.Post("/logout", handlers.NewPostLogoutHandler(handlers.PostLogoutHandlerParams{
			SessionCookie: s.sessionCookie,
		}).ServeHTTP)

		// Todos Handler
		r.Route("/todos", func(r chi.Router) {

			todoHandler := handlers.NewTodoHandler(handlers.TodoHandlerParams{
				Users: s.users,
				Todos: s.todos,
			})

			r.Get("/", todoHandler.List)
			r.Post("/", todoHandler.Create)
			r.Delete("/{id}", todoHandler.Delete)
			r.Delete("/all", todoHandler.DeleteAll)
			r.Delete("/item/{todoId}/{itemId}", todoHandler.DeleteItem)
			r.Get("/item/edit/{id}", todoHandler.EditItem)
			r.Post("/item/edit/content/{id}", todoHandler.UpdateItemContent)
			r.Post("/item/toggle/{id}", todoHandler.ToggleItemCheck)
			r.Post("/item/{id}", todoHandler.AddItem)
		})
	})
} */

func (s *ZionServer) EsatblishRoutesV2() {
	mux := s.router

	authentication := auth.NewAuthMiddleware(auth.AuthMiddlewareParams{
		Sessions:          s.sessions,
		SessionCookieName: s.sessionCookie,
	})

	middleware := func(h http.Handler) http.Handler {
		return authentication.AddUserToContext(
			middleware.Logger(
				middleware.Recoverer(h),
			),
		)
	}

	// Index Handler
	mux.HandleFunc("/", middleware(handlers.NewIndex()).ServeHTTP)

	// Login Handler
	mux.HandleFunc("/login", middleware(handlers.NewGetLoginHandler()).ServeHTTP)
	mux.HandleFunc("POST /login", middleware(
		handlers.NewPostLoginHandler(
			handlers.PostLoginHandlerParameters{
				Users:         s.users,
				Sessions:      s.sessions,
				PasswordHash:  s.hash,
				SessionCookie: s.sessionCookie,
			})).ServeHTTP)

	// Register Handler
	mux.HandleFunc("/register", middleware(handlers.NewGetRegisterHandler()).ServeHTTP)
	mux.HandleFunc("POST /register", middleware(
		handlers.NewPostRegisterHandler(
			handlers.PostRegisterHandlerParameters{
				Users: s.users,
			})).ServeHTTP)

	// Logout Handler
	mux.HandleFunc("POST /logout", middleware(
		handlers.NewPostLogoutHandler(
			handlers.PostLogoutHandlerParams{
				SessionCookie: s.sessionCookie,
			})).ServeHTTP)

	// Todos Handler
	todoHandler := handlers.NewTodoHandler(handlers.TodoHandlerParams{
		Users: s.users,
		Todos: s.todos,
	})

	todoGroupHandler := http.NewServeMux()
	todoGroupHandler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		middleware(http.HandlerFunc(todoHandler.List)).ServeHTTP(w, r)
	})

	mux.Handle("/todos/", http.StripPrefix("/todos", todoGroupHandler))
}
