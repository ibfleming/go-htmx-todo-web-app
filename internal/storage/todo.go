package storage

import (
	"zion/internal/storage/db"

	"gorm.io/gorm"
)

type TodoStorage struct {
	db *gorm.DB
}

type TodoStorageParams struct {
	DB *gorm.DB
}

func NewTodoStorage(params TodoStorageParams) *TodoStorage {
	return &TodoStorage{
		db: params.DB,
	}
}

func (h *TodoStorage) CreateTodo(todo db.Todo) (*db.Todo, error) {
	if err := h.db.Create(&todo).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

func (h *TodoStorage) AddTodoItemToTodo(todoID uint, item *db.TodoItem) (*db.TodoItem, error) {
	return nil, nil
}

func (h *TodoStorage) DeleteTodo(todoID string, userID uint) error {
	err := h.db.Where("id = ? AND user_id = ?", todoID, userID).Delete(&db.Todo{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (h *TodoStorage) DeleteAllTodos(userID uint) error {
	err := h.db.Where("user_id = ?", userID).Delete(&db.Todo{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (h *TodoStorage) DeleteTodoItem(todoID, itemID uint) error {
	return nil
}

func (h *TodoStorage) GetTodosByUserID(userID uint) ([]*db.Todo, error) {
	var todos []*db.Todo
	err := h.db.Where("user_id = ?", userID).Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (h *TodoStorage) GetTodoByTodoID(todoID uint) (*db.Todo, error) {
	return nil, nil
}

func (h *TodoStorage) UpdateTodo(todoID uint, title, description string) error {
	return nil
}

func (h *TodoStorage) UpdateTodoItem(todoID, itemID uint, checked bool, content string) error {
	return nil
}

func (h *TodoStorage) ListTodoItems(todoID uint) ([]*db.TodoItem, error) {
	return nil, nil
}
