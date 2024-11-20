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

func (h *TodoStorage) AddTodoItemToTodo(item *schema.TodoItem) (*schema.TodoItem, error) {
	if err := h.db.Create(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
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

func (h *TodoStorage) DeleteTodoItemByID(itemID string) error {
	err := h.db.Where("id = ?", itemID).Delete(&schema.TodoItem{}).Error
	if err != nil {
		return err
	}
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

func (h *TodoStorage) GetTodoItemByID(itemID string) (*schema.TodoItem, error) {
	var item *schema.TodoItem
	err := h.db.Where("id = ?", itemID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *TodoStorage) GetTodoByTodoID(todoID uint) (*schema.Todo, error) {
	return nil, nil
}

func (h *TodoStorage) UpdateTodo(todoID uint, title, description string) error {
	return nil
}

func (h *TodoStorage) UpdateTodoItemContent(itemID string, content string) (*schema.TodoItem, error) {
	item, err := h.GetTodoItemByID(itemID)
	if err != nil {
		return nil, err
	}
	item.Content = content
	err = h.db.Save(&item).Error
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *TodoStorage) UpdateTodoItemChecked(itemID string, checked bool) (*schema.TodoItem, error) {
	item, err := h.GetTodoItemByID(itemID)
	if err != nil {
		return nil, err
	}
	item.Checked = checked
	err = h.db.Save(&item).Error
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *TodoStorage) ListTodoItems(todoID uint) ([]*schema.TodoItem, error) {
	return nil, nil
}

func (h *TodoStorage) GetTodoItemLenthByID(todoID uint) (int, error) {
	var todo *schema.Todo
	err := h.db.Preload("Items").Where("id = ?", todoID).First(&todo).Error
	if err != nil {
		return 0, err
	}
	return len(todo.Items), nil
}
