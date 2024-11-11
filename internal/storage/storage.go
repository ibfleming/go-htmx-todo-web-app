package storage

import "zion/internal/storage/db"

type UserStorageInterface interface {
	CreateUser(email string, password string) error
	GetUser(email string) (*db.User, error)
}

type SessionStorageInterface interface {
	CreateSession(session *db.Session) (*db.Session, error)
	GetUserFromSession(sessionID, userID string) (*db.User, error)
	DeleteSession(sessionID string) error
}

type TodoStorageInterface interface {
	CreateTodo(todo *db.Todo) (*db.Todo, error)
	CreateTodoItem(todoItem *db.TodoItem) (*db.TodoItem, error)
	DeleteTodo(todoId string) error
	DeleteTodoItem(todoItemId string) error
	DeleteChecklistItem(todoId, itemId uint) error
	GetAllTodos(userId string) ([]db.Todo, error)
	GetTodoItems(todoId string) ([]db.TodoItem, error)
}
