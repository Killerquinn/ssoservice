package storage

import "errors"

var (
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrAppNotFound        = errors.New("app not found")
	ErrDoesntAllowed      = errors.New("doesnt allowed for this role")
	ErrInvalidCredentials = errors.New("error invalid credentials")
)
