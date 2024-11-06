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
	AddTodoItemToTodo(todoID uint, item *db.TodoItem) (*db.TodoItem, error)
	DeleteTodo(todoID uint) error
	DeleteChecklistItem(todoID, itemID uint) error
	GetTodos(userID uint) ([]db.Todo, error)
}
