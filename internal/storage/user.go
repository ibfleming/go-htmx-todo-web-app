package storage

import (
	zerr "zion/internal/errors"
	"zion/internal/hash"
	schema "zion/internal/storage/schema"

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
	return s.db.Create(&schema.User{
		Email:    email,
		Password: hashedPassword,
	}).Error
}

func (s *UserStorage) GetUser(email string) (*schema.User, error) {
	var user schema.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, zerr.ErrUserNotFound
	}
	return &user, nil
}

func (s *UserStorage) GetUserByID(userID uint) (*schema.User, error) {
	var user schema.User
	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, zerr.ErrUserNotFound
	}
	return &user, nil
}

func (s *UserStorage) UpdateUser(userID uint, email, password string) error {
	return nil
}

func (s *UserStorage) DeleteUser(userID uint) error {
	return nil
}

func (s *UserStorage) UserExists(email string) (bool, error) {
	return false, nil
}
