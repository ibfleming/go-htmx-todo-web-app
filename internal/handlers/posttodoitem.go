package handlers

import (
	"net/http"
	"strconv"
	"zion/internal/storage"
	"zion/internal/storage/db"
	"zion/templates"

	"github.com/go-chi/chi/v5"
)

type PostTodoItemHandler struct {
	todos storage.TodoStorageInterface
}

type PostTodoItemHandlerParams struct {
	Todos storage.TodoStorageInterface
}

func NewPostTodoItemHandler(params PostTodoItemHandlerParams) *PostTodoItemHandler {
	return &PostTodoItemHandler{
		todos: params.Todos,
	}
}

func (h *PostTodoItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	todoId := chi.URLParam(r, "todoId")

	tempTodoId, err := strconv.ParseUint(todoId, 10, 32)
	if err != nil {
		http.Error(w, "invalid todoId", http.StatusBadRequest)
		return
	}

	_, err = h.todos.CreateTodoItem(&db.TodoItem{
		TodoID:      uint(tempTodoId),
		Description: "",
		Completed:   false,
	})
	if err != nil {
		http.Error(w, "error creating todo item", http.StatusInternalServerError)
		return
	}

	todoItems, err := h.todos.GetTodoItems(todoId)
	if err != nil {
		http.Error(w, "error fetching todo items", http.StatusInternalServerError)
		return
	}

	err = templates.TodoItemList(todoItems).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}
