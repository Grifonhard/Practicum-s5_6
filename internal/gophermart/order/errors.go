package order

import "errors"

var (
	ErrLuhnFail = errors.New("order number didn 't pass Luhn algoritm test")
	ErrOrderNotReady = errors.New("reward for the order has not yet been calculated")
	ErrOrderInvalid = errors.New("order is invalid")
)