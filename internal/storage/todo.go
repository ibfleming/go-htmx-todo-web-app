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

func (s *TodoStorage) CreateTodo(todo *db.Todo) (*db.Todo, error) {
	result := s.db.Create(todo)
	if result.Error != nil {
		return nil, result.Error
	}
	return todo, nil
}

func (s *TodoStorage) AddTodoItemToTodo(todoID uint, item *db.TodoItem) (*db.TodoItem, error) {
	// Find the parent TODO by ID
	var todo db.Todo
	if err := s.db.First(&todo, todoID).Error; err != nil {
		return nil, err
	}
	// TODO is valid, now add TodoItem to it
	item.TodoID = todoID
	todo.Items = append(todo.Items, *item)
	// Save the TodoItem to the database
	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}
	// Save the updated Todo to the database
	if err := s.db.Save(&todo).Error; err != nil {
		return nil, err
	}
	// Return the created TodoItem
	return item, nil
}

func (s *TodoStorage) DeleteTodo(todoID string) error {
	var todo db.Todo
	if err := s.db.First(&todo, todoID).Error; err != nil {
		return err
	}
	if err := s.db.Delete(&todo).Error; err != nil {
		return err
	}
	return nil
}

func (s *TodoStorage) DeleteChecklistItem(todoID, itemID uint) error {
	var item db.TodoItem
	if err := s.db.Where("id = ? AND todo_id = ?", itemID, todoID).First(&item).Error; err != nil {
		return err
	}
	if err := s.db.Delete(&item).Error; err != nil {
		return err
	}
	return nil
}

func (s *TodoStorage) GetTodos(userID string) ([]db.Todo, error) {
	var todos []db.Todo
	if err := s.db.Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}
