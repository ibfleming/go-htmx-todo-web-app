package errors

import "errors"

var (
	ErrFailedToConnectToDB  = errors.New("failed to connect to the database")
	ErrDropTables           = errors.New("failed to drop tables")
	ErrCreateTables         = errors.New("failed to create tables")
	ErrInvalidHash          = errors.New("encoded hash is in the incorrect format")
	ErrIncompatibleVersion  = errors.New("incompatible version of Argon2")
	ErrUserNotFound         = errors.New("user not found")
	ErrHashPasswordFailed   = errors.New("failed to hash password")
	ErrInvalidCookie        = errors.New("invalid cookie")
	ErrFailedToDecodeString = errors.New("failed to decode string")
	ErrInvalidCookieFormat  = errors.New("invalid cookie format")
	ErrFailedToGetUserFromSession = errors.New("failed to get user from session")
)
