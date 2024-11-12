package storage

import (
	"errors"
	"log"
	"strconv"
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

func (s *TodoStorage) DeleteTodoItem(todoId, itemId string) (error, bool) {
	todoIdUInt, err := strconv.ParseUint(todoId, 10, 32)
	if err != nil {
		log.Printf("%v\n", err)
		return err, false
	}

	itemIdUInt, err := strconv.ParseUint(itemId, 10, 32)
	if err != nil {
		return err, false
	}

	var todo db.Todo
	if err := s.db.Preload("Items").First(&todo, todoIdUInt).Error; err != nil {
		return err, false
	}

	var itemDelete *db.TodoItem
	for i, item := range todo.Items {
		if item.ID == uint(itemIdUInt) {
			itemDelete = &todo.Items[i]
			break
		}
	}

	if itemDelete == nil {
		return errors.New("item not found"), false
	}

	if err := s.db.Delete(itemDelete).Error; err != nil {
		return err, false
	}

	if err := s.db.Preload("Items").First(&todo, todoIdUInt).Error; err != nil {
		return err, false
	}

	if len(todo.Items) == 0 {
		return nil, true
	}

	return nil, false
}

func (s *TodoStorage) GetAllTodosForUser(userId uint) ([]db.Todo, error) {
	var todos []db.Todo
	if err := s.db.
		Where(&db.Todo{UserID: userId}).
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
