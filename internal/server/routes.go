package server

import (
	"net/http"
	"zion/internal/handlers"
	"zion/internal/middleware/auth"

	"github.com/go-chi/chi/v5/middleware"
)

func (s *ZionServer) EsatblishRoutes() {
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
		switch r.Method {
		case http.MethodGet:
			middleware(http.HandlerFunc(todoHandler.List)).ServeHTTP(w, r)
		case http.MethodPost:
			middleware(http.HandlerFunc(todoHandler.Create)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	todoGroupHandler.HandleFunc("DELETE /{id}", func(w http.ResponseWriter, r *http.Request) {
		middleware(http.HandlerFunc(todoHandler.Delete)).ServeHTTP(w, r)
	})
	todoGroupHandler.HandleFunc("DELETE /all", func(w http.ResponseWriter, r *http.Request) {
		middleware(http.HandlerFunc(todoHandler.DeleteAll)).ServeHTTP(w, r)
	})
	todoGroupHandler.HandleFunc("DELETE /item/{todoId}/{itemId}", func(w http.ResponseWriter, r *http.Request) {
		middleware(http.HandlerFunc(todoHandler.DeleteItem)).ServeHTTP(w, r)
	})

	todoGroupHandler.HandleFunc("GET /item/edit/{id}", func(w http.ResponseWriter, r *http.Request) {
		middleware(http.HandlerFunc(todoHandler.EditItem)).ServeHTTP(w, r)
	})

	todoGroupHandler.HandleFunc("POST /item/edit/content/{id}", func(w http.ResponseWriter, r *http.Request) {
		middleware(http.HandlerFunc(todoHandler.UpdateItemContent)).ServeHTTP(w, r)
	})

	todoGroupHandler.HandleFunc("POST /item/toggle/{id}", func(w http.ResponseWriter, r *http.Request) {
		middleware(http.HandlerFunc(todoHandler.ToggleItemCheck)).ServeHTTP(w, r)
	})

	todoGroupHandler.HandleFunc("POST /item/{id}", func(w http.ResponseWriter, r *http.Request) {
		middleware(http.HandlerFunc(todoHandler.AddItem)).ServeHTTP(w, r)
	})

	mux.Handle("/todos/", http.StripPrefix("/todos", todoGroupHandler))
}
