package handlers

import (
	"net/http"
	"zion/internal/storage"
	"zion/templates"

	"github.com/go-chi/chi/v5"
)

type GetTodoHandler struct {
	todos storage.TodoStorageInterface
}

type GetTodoHandlerParameters struct {
	Todos storage.TodoStorageInterface
}

func NewGetTodoHandler(params GetTodoHandlerParameters) *GetTodoHandler {
	return &GetTodoHandler{
		todos: params.Todos,
	}
}

func (h *GetTodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	todos, err := h.todos.GetTodos(userId)

	if err != nil {
		http.Error(w, "error fetching todos", http.StatusInternalServerError)
		return
	}

	if len(todos) == 0 {
		err = templates.EmptyTodoList().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
			return
		}
	}

	err = templates.TodosList(todos).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}
