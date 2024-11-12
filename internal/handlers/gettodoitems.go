package handlers

import (
	"net/http"
	"zion/internal/storage"
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
}
