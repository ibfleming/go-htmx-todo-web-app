package handlers

import (
	"net/http"
	"zion/internal/storage"

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
}
