package storage

import (
	"zion/internal/storage/schema"
)

type UserStorageInterface interface {
	CreateUser(email, password string) error
	GetUser(email string) (*schema.User, error)
	GetUserByID(userID uint) (*schema.User, error)
	UpdateUser(userID uint, email, password string) error
	DeleteUser(userID uint) error
	UserExists(email string) (bool, error)
}

type SessionStorageInterface interface {
	CreateSession(session *schema.Session) (*schema.Session, error)
	GetUserFromSession(sessionID, userID string) (*schema.User, error)
	DeleteSession(sessionID string) error
}

type TodoStorageInterface interface {
	CreateTodo(todo schema.Todo) (*schema.Todo, error)
	AddTodoItemToTodo(item *schema.TodoItem) (*schema.TodoItem, error)
	DeleteTodo(todoID string, userID uint) error
	DeleteAllTodos(userID uint) error
	DeleteTodoItemByID(itemID string) error
	GetTodosByUserID(userID uint) ([]*schema.Todo, error)
	GetTodoByTodoID(todoID uint) (*schema.Todo, error)
	GetTodoItemByID(itemID string) (*schema.TodoItem, error)
	GetTodoItemLenthByID(todoID uint) (int, error)
	UpdateTodo(todoID uint, title, description string) error
	UpdateTodoItemContent(itemID string, content string) (*schema.TodoItem, error)
	UpdateTodoItemChecked(itemID string, checked bool) error
	ListTodoItems(todoID uint) ([]*schema.TodoItem, error)
}
