package storage

import (
	"zion/internal/storage/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionStorage struct {
	db *gorm.DB
}

type SessionStorageParameters struct {
	DB *gorm.DB
}

func NewSessionStorage(params SessionStorageParameters) *SessionStorage {
	return &SessionStorage{
		db: params.DB,
	}
}

func (s *SessionStorage) CreateSession(session *db.Session) (*db.Session, error) {
	session.SessionID = uuid.New().String()
	result := s.db.Create(session)
	if result.Error != nil {
		return nil, result.Error
	}
	return session, nil
}

func (s *SessionStorage) GetUserFromSession(sessionID, userID string) (*db.User, error) {
	var session db.Session
	err := s.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Email")
	}).Where("session_id = ? AND user_id = ?", sessionID, userID).First(&session).Error
	if err != nil {
		return nil, err
	}
	if session.User.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &session.User, nil
}
