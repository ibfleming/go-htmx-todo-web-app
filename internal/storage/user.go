package storage

import (
	"zion/internal/storage/db"

	"gorm.io/gorm"
)

type UserStorage struct {
	db           *gorm.DB
	passwordHash string
}

type UserStorageParameters struct {
	DB           *gorm.DB
	PasswordHash string
}

func NewUserStorage(params UserStorageParameters) *UserStorage {
	return &UserStorage{
		db:           params.DB,
		passwordHash: params.PasswordHash,
	}
}

func (s *UserStorage) CreateUser(email, password string) error {
	return s.db.Create(&db.User{
		Email:    email,
		Password: password,
	}).Error
}

func (s *UserStorage) GetUser(email string) (*db.User, error) {
	var user db.User

	err := s.db.Where(&db.User{Email: email}).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, err
}
