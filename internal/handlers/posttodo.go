package handlers

import (
	"fmt"
	"log"
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
	user := auth.GetUser(r.Context())
	title := r.FormValue("title")
	description := r.FormValue("description")

	log.Printf("User: %s, Title: %s, Description: %s", user.Email, title, description)

	_, err := h.todos.CreateTodo(&db.Todo{
		Title:       title,
		Description: description,
		UserID:      user.ID,
		User:        *user,
	})
	if err != nil {
		http.Error(w, "error creating todo", http.StatusInternalServerError)
		return
	}

	todos, err := h.todos.GetAllTodos(fmt.Sprintf("%d", user.ID))
	if err != nil {
		http.Error(w, "error fetching todos", http.StatusInternalServerError)
		return
	}

	err = templates.TodoList(todos).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}
