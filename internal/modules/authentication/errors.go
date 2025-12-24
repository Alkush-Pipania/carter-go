package authentication

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUnauthorized       = errors.New("unauthorized")
)
