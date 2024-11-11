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
	if err := s.db.Create(todo).Error; err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *TodoStorage) CreateTodoItem(todoItem *db.TodoItem) (*db.TodoItem, error) {
	if err := s.db.Create(todoItem).Error; err != nil {
		return nil, err
	}
	return todoItem, nil
}

func (s *TodoStorage) DeleteTodo(todoId string) error {
	if err := s.db.Delete(&db.Todo{}, todoId).Error; err != nil {
		return err
	}
	return nil
}

func (s *TodoStorage) DeleteTodoItem(todoItemId string) error {
	return nil
}

func (s *TodoStorage) DeleteChecklistItem(todoId, itemId uint) error {
	return nil
}

func (s *TodoStorage) GetAllTodos(userId string) ([]db.Todo, error) {
	var todos []db.Todo
	if err := s.db.
		Where("user_id = ?", userId).
		Preload("Items").
		Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (s *TodoStorage) GetTodoItems(todoId string) ([]db.TodoItem, error) {
	var todoItems []db.TodoItem
	if err := s.db.
		Where("todo_id = ?", todoId).
		Find(&todoItems).Error; err != nil {
		return nil, err
	}
	return todoItems, nil
}
