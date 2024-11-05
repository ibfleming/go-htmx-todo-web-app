package storage

import (
	zerr "zion/internal/errors"
	"zion/internal/hash"
	"zion/internal/storage/db"

	"gorm.io/gorm"
)

type UserStorage struct {
	db           *gorm.DB
	passwordHash *hash.PasswordHash
}

type UserStorageParams struct {
	DB           *gorm.DB
	PasswordHash *hash.PasswordHash
}

func NewUserStorage(params UserStorageParams) *UserStorage {
	return &UserStorage{
		db:           params.DB,
		passwordHash: params.PasswordHash,
	}
}

func (s *UserStorage) CreateUser(email, password string) error {
	hashedPassword, err := s.passwordHash.GenerateFromPassword(password)
	if err != nil {
		return zerr.ErrHashPasswordFailed
	}
	return s.db.Create(&db.User{
		Email:    email,
		Password: hashedPassword,
	}).Error
}

func (s *UserStorage) GetUser(email string) (*db.User, error) {
	var user db.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, zerr.ErrUserNotFound
	}
	return &user, nil
}
