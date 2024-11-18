package storage

import "errors"

var (
	ErrUserExist    = errors.New("such user already exists")
	ErrUserNotExist = errors.New("user not found")
)
