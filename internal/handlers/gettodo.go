package handlers

import (
	"net/http"
	"zion/internal/storage"
)

type GetTodoHandler struct {
	todos storage.TodoStorageInterface
}

type GetTodoHandlerParams struct {
	Todos storage.TodoStorageInterface
}

func NewGetTodoHandler(params GetTodoHandlerParams) *GetTodoHandler {
	return &GetTodoHandler{
		todos: params.Todos,
	}
}

func (h *GetTodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
