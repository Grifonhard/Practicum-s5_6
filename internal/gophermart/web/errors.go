package web

import "errors"

var (
	ErrUserNotFoundCtx = errors.New("user id not found in context")
	ErrUserIDWrongType = errors.New("user id type assertion failed")
)