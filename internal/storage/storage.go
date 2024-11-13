package storage

import "zion/internal/storage/db"

type UserStorageInterface interface {
	CreateUser(email, password string) error
	GetUser(email string) (*db.User, error)
	GetUserByID(userID uint) (*db.User, error)
	UpdateUser(userID uint, email, password string) error
	DeleteUser(userID uint) error
	UserExists(email string) (bool, error)
}

type SessionStorageInterface interface {
	CreateSession(session *db.Session) (*db.Session, error)
	GetUserFromSession(sessionID, userID string) (*db.User, error)
	DeleteSession(sessionID string) error
}

type TodoStorageInterface interface {
	CreateTodo(todo db.Todo) (*db.Todo, error)
	AddTodoItemToTodo(todoID uint, item *db.TodoItem) (*db.TodoItem, error)
	DeleteTodo(todoID string, userID uint) error
	DeleteAllTodos(userID uint) error
	DeleteTodoItem(todoID, itemID uint) error
	GetTodosByUserID(userID uint) ([]*db.Todo, error)
	GetTodoByTodoID(todoID uint) (*db.Todo, error)
	UpdateTodo(todoID uint, title, description string) error
	UpdateTodoItem(todoID, itemID uint, checked bool, content string) error
	ListTodoItems(todoID uint) ([]*db.TodoItem, error)
}
