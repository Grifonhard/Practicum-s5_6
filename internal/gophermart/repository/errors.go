package psql

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrOrderNotFound  = errors.New("order not found")
	ErrOrdersNotFound = errors.New("orders not found")
	ErrTransNotFound  = errors.New("transactions not found")
	ErrUserExist      = errors.New("user exist already")
	ErrOrderExist     = errors.New("order exist already")
)
