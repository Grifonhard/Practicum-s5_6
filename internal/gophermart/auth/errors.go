package auth

import "errors"

var (
	ErrWrongPassword = errors.New("wrong password")
	ErrInvalidToken  = errors.New("invalid token")
)
