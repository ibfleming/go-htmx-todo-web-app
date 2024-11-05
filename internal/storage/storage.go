package storage

import "zion/internal/storage/db"

type UserStorageInterface interface {
	CreateUser(email string, password string) error
	GetUser(email string) (*db.User, error)
}

type SessionStorageInterface interface {
	CreateSession(session *db.Session) (*db.Session, error)
	GetUserFromSession(sessionID, userID string) (*db.User, error)
	DeleteSession(sessionID string) error
}
