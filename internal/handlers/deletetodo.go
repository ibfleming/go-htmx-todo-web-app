package handlers

import (
	"net/http"
	"strconv"
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
	todoId := chi.URLParam(r, "id")
	user := auth.GetUser(r.Context())
	userId := strconv.FormatUint(uint64(user.ID), 10)

	err := h.todos.DeleteTodo(todoId)
	if err != nil {
		http.Error(w, "error deleting todo", http.StatusInternalServerError)
		return
	}

	todo, _ := h.todos.GetTodos(userId)
	if len(todo) == 0 {
		err = templates.EmptyTodoList().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
			return
		}
	}
}
