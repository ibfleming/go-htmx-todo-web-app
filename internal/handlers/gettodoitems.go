package handlers

import (
	"net/http"
	"zion/internal/storage"
	"zion/templates"

	"github.com/go-chi/chi/v5"
)

type GetTodoItemsHandler struct {
	todos storage.TodoStorageInterface
}

type GetTodoItemsHandlerParams struct {
	Todos storage.TodoStorageInterface
}

func NewGetTodoItemsHandler(params GetTodoItemsHandlerParams) *GetTodoItemsHandler {
	return &GetTodoItemsHandler{
		todos: params.Todos,
	}
}

func (h *GetTodoItemsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	todoId := chi.URLParam(r, "todoId")
	todoItems, err := h.todos.GetTodoItems(todoId)

	if err != nil {
		http.Error(w, "error fetching todo items", http.StatusInternalServerError)
		return
	}

	if len(todoItems) == 0 {
		err = templates.EmptyTodoItemListt().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
			return
		}
	}

	err = templates.TodoItemList(todoItems).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}

}
