package auth

import "errors"

var (
	ErrWrongPassword = errors.New("wrong password")
	ErrInvalidToken  = errors.New("invalid token")
	ErrUserExist    = errors.New("such user already exists")
	ErrUserNotExist = errors.New("user not found")
)
