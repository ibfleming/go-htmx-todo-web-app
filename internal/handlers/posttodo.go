package handlers

import (
	"net/http"
	"zion/internal/middleware/auth"
	"zion/internal/storage"
	"zion/internal/storage/db"
	"zion/templates"
)

type PostTodoHandler struct {
	todos storage.TodoStorageInterface
}

type PostTodoHandlerParams struct {
	Todos storage.TodoStorageInterface
}

func NewPostTodoHandler(params PostTodoHandlerParams) *PostTodoHandler {
	return &PostTodoHandler{
		todos: params.Todos,
	}
}

func (h *PostTodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// (1) Get user, title, and description
	user := auth.GetUser(r.Context())
	title := r.FormValue("title")
	description := r.FormValue("description")

	// (2) Create the todo (addsto database)
	todo, err := h.todos.CreateTodo(&db.Todo{
		Title:       title,
		Description: description,
		UserID:      user.ID,
		User:        *user,
	})
	if err != nil {
		http.Error(w, "error creating todo", http.StatusInternalServerError)
		return
	}

	// (3) Return the html
	err = templates.SingleTodo(*todo).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}
