package handlers

import (
	"net/http"
	"zion/internal/middleware/auth"
	"zion/internal/storage"
	"zion/templates"

	"github.com/go-chi/chi/v5"
)

type DeleteTodoHandler struct {
	todos storage.TodoStorageInterface
}

type DeleteTodoHandlerParams struct {
	Todos storage.TodoStorageInterface
}

func NewDeleteTodoHandler(params DeleteTodoHandlerParams) *DeleteTodoHandler {
	return &DeleteTodoHandler{
		todos: params.Todos,
	}
}

func (h *DeleteTodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.todos.DeleteTodo(chi.URLParam(r, "todoId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todos, _ := h.todos.GetAllTodosForUser(auth.GetUserID(r.Context()))
	if len(todos) == 0 {
		err := templates.EmptyTodoList().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
			return
		}
	}
}
