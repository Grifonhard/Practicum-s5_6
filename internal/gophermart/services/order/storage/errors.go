package storage

import "errors"

var (
	ErrUserExist    = errors.New("such user already exists")
	ErrUserNotExist = errors.New("user not found")
	ErrOrderExistThis = errors.New("order whith this id of this user is exist")
	ErrOrderExistAnother = errors.New("order whith this id of another user is exist")
)
