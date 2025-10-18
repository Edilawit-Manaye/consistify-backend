package domain

import "errors"
var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrPasswordsDoNotMatch   = errors.New("passwords do not match")
	ErrInvalidToken          = errors.New("invalid token")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
	ErrConsistencyNotFound   = errors.New("consistency record not found")
	ErrPlatformNotLinked     = errors.New("platform not linked for user")
	ErrExternalAPIFailed     = errors.New("external platform API failed")
	ErrProcessingConsistency = errors.New("error processing consistency data")
	ErrInvalidNotificationTime = errors.New("invalid notification time format, expected HH:MM")
)



