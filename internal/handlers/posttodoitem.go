package handlers

import (
	"log"
	"net/http"
	"zion/internal/storage"
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
	/*
		/todo/{todoId}/item/create
			1. get the todo parent
			2. if valid create the todo itemm
			3. add todo item to database
			4. associate the item to the parent todo
			5. return the html li element
	*/
	log.Printf("PostTodoItemHandler: %v", r.URL.Path)
}
