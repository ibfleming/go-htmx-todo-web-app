package storage

import (
	zerr "zion/internal/errors"
	schema "zion/internal/storage/schema"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionStorage struct {
	db *gorm.DB
}

type SessionStorageParams struct {
	DB *gorm.DB
}

func NewSessionStorage(params SessionStorageParams) *SessionStorage {
	return &SessionStorage{
		db: params.DB,
	}
}

func (s *SessionStorage) CreateSession(session *schema.Session) (*schema.Session, error) {
	session.SessionID = uuid.New().String()
	result := s.db.Create(session)
	if result.Error != nil {
		return nil, result.Error
	}
	return session, nil
}

func (s *SessionStorage) GetUserFromSession(sessionID, userID string) (*schema.User, error) {
	var session schema.Session
	err := s.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Email")
	}).Where("session_id = ? AND user_id = ?", sessionID, userID).First(&session).Error
	if err != nil {
		return nil, err
	}
	if session.User.ID == 0 {
		return nil, zerr.ErrUserNotFound
	}
	return &session.User, nil
}

func (s *SessionStorage) DeleteSession(sessionID string) error {
	return s.db.Where("session_id = ?", sessionID).Delete(&schema.Session{}).Error
}
