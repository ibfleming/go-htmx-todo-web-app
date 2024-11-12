package server

import (
	"zion/internal/handlers"
	"zion/internal/middleware/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *ZionServer) EstablishRoutes() {

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
		r.NotFound(handlers.NewNotFoundHandler().ServeHTTP)

		// Home Handler
		r.Get("/", handlers.NewGetIndexHandler(handlers.IndexParams{
			Todos: s.todos,
		}).ServeHTTP)

		// Login Handler
		r.Get("/login", handlers.NewGetLoginHandler().ServeHTTP)
		r.Post("/login", handlers.NewPostLoginHandler(handlers.PostLoginHandlerParams{
			Users:         s.users,
			Sessions:      s.sessions,
			PasswordHash:  s.hash,
			SessionCookie: s.sessionCookie,
		}).ServeHTTP)

		// Register Handler
		r.Get("/register", handlers.NewGetRegisterHandler().ServeHTTP)
		r.Post("/register", handlers.NewPostRegisterHandler(handlers.PostRegisterHandlerParams{
			Users: s.users,
		}).ServeHTTP)

		// Logout Handler
		r.Post("/logout", handlers.NewPostLogoutHandler(handlers.PostLogoutHandlerParams{
			SessionCookie: s.sessionCookie,
		}).ServeHTTP)

		// Get All Todos Handler
		// r.Get("/todo", handlers.NewGetTodoHandler(handlers.GetTodoHandlerParams{
		// 	Todos: s.todos,
		// }).ServeHTTP)

		// Add Todo Handler
		r.Post("/todo/create", handlers.NewPostTodoHandler(handlers.PostTodoHandlerParams{
			Todos: s.todos,
		}).ServeHTTP)

		// Delete Todo Handler
		r.Delete("/todo/{todoId}", handlers.NewDeleteTodoHandler(handlers.DeleteTodoHandlerParams{
			Todos: s.todos,
		}).ServeHTTP)

		r.Post("/todo/{todoId}/item/create", handlers.NewPostTodoItemHandler(handlers.PostTodoItemHandlerParams{
			Todos: s.todos,
		}).ServeHTTP)

		r.Delete("/todo/{todoId}/item/{itemId}", handlers.NewDeleteTodoItemHandler(handlers.DeleteTodoItemHandlerParams{
			Todos: s.todos,
		}).ServeHTTP)

		/*
			// Get All Todo Checklist Items
			r.Get("/{userId}/todo/{todoId}/item", handlers.NewGetTodoItemsHandler(handlers.GetTodoItemsHandlerParams{
				Todos: s.todos,
			}).ServeHTT
			// Add Todo Checklist Item
			r.Post("/{userId}/todo/{todoId}/item", handlers.NewPostTodoItemHandler(handlers.PostTodoItemHandlerParams{
				Todos: s.todos,
			}).ServeHTTP)
		*/
	})
}
