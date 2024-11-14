package storage

import (
	schema "zion/internal/storage/schema"

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

func (h *TodoStorage) CreateTodo(todo schema.Todo) (*schema.Todo, error) {
	if err := h.db.Create(&todo).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

func (h *TodoStorage) AddTodoItemToTodo(todoID uint, item *schema.TodoItem) (*schema.TodoItem, error) {
	return nil, nil
}

func (h *TodoStorage) DeleteTodo(todoID string, userID uint) error {
	err := h.db.Where("id = ? AND user_id = ?", todoID, userID).Delete(&schema.Todo{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (h *TodoStorage) DeleteAllTodos(userID uint) error {
	err := h.db.Where("user_id = ?", userID).Delete(&schema.Todo{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (h *TodoStorage) DeleteTodoItem(todoID, itemID uint) error {
	return nil
}

func (h *TodoStorage) GetTodosByUserID(userID uint) ([]*schema.Todo, error) {
	var todos []*schema.Todo
	err := h.db.Preload("Items").Where("user_id = ?", userID).Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (h *TodoStorage) GetTodoByTodoID(todoID uint) (*schema.Todo, error) {
	return nil, nil
}

func (h *TodoStorage) UpdateTodo(todoID uint, title, description string) error {
	return nil
}

func (h *TodoStorage) UpdateTodoItem(todoID, itemID uint, checked bool, content string) error {
	return nil
}

func (h *TodoStorage) ListTodoItems(todoID uint) ([]*schema.TodoItem, error) {
	return nil, nil
}
