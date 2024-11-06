package errors

import (
	"errors"
	"regexp"
)

var (
	ErrFailedToConnectToDB        = errors.New("failed to connect to the database")
	ErrDropTables                 = errors.New("failed to drop tables")
	ErrCreateTables               = errors.New("failed to create tables")
	ErrInvalidHash                = errors.New("encoded hash is in the incorrect format")
	ErrIncompatibleVersion        = errors.New("incompatible version of Argon2")
	ErrUserNotFound               = errors.New("user not found")
	ErrHashPasswordFailed         = errors.New("failed to hash password")
	ErrInvalidCookie              = errors.New("no session cookie found")
	ErrFailedToDecodeString       = errors.New("failed to decode string")
	ErrInvalidCookieFormat        = errors.New("invalid cookie format")
	ErrFailedToGetUserFromSession = errors.New("failed to get user from session")
	ErrEmailAlreadyInUse          = errors.New("email already in use")
	ErrSQLQueryFailed             = errors.New("bad request")
	ErrPasswordIncorrect          = errors.New("password is incorrect")
	ErrInvalidEmailOrPassword     = errors.New("invalid email or password")
)

func ExtractSQLStateErrorCode(errMsg string) string {
	re := regexp.MustCompile(`SQLSTATE (\d{5})`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func SQLErrorMessage(err error) string {
	statusCode := ExtractSQLStateErrorCode(err.Error())
	switch statusCode {
	case "23505":
		return ErrEmailAlreadyInUse.Error()
	default:
		return ErrSQLQueryFailed.Error()
	}
}
