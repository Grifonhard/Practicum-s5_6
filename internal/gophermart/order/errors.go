package order

import "errors"

var (
	ErrLuhnFail = errors.New("order number didn 't pass Luhn algoritm test")
)