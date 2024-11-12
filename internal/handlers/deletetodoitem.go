package handlers

import (
	"log"
	"net/http"
	"zion/internal/storage"
	"zion/templates"

	"github.com/go-chi/chi/v5"
)

type DeleteTodoItemHandler struct {
	todos storage.TodoStorageInterface
}

type DeleteTodoItemHandlerParams struct {
	Todos storage.TodoStorageInterface
}

func NewDeleteTodoItemHandler(params DeleteTodoItemHandlerParams) *DeleteTodoItemHandler {
	return &DeleteTodoItemHandler{
		todos: params.Todos,
	}
}

func (h *DeleteTodoItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// "/todo/{todoId}/item/{itemId}"
	todoId := chi.URLParam(r, "todoId")
	itemId := chi.URLParam(r, "itemId")

	err, empty := h.todos.DeleteTodoItem(todoId, itemId)
	if err != nil {
		log.Print(err)
		http.Error(w, "error deleting todo item", http.StatusInternalServerError)
		return
	}

	if empty {
		err = templates.EmptyTodoItems().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
			return
		}
	}
}
