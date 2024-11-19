package order

import "errors"

var (
	ErrLuhnFail         = errors.New("order number didn 't pass Luhn algoritm test")
	ErrOrderNotReady    = errors.New("reward for the order has not yet been calculated")
	ErrOrderInvalid     = errors.New("order is invalid")
	ErrNotEnoughBalance = errors.New("not enough balanse for withdraw")
	ErrAlreadyDebited   = errors.New("points for this order have already been debited")
	ErrTooMuchTransact  = errors.New("too much transactions by one order")	
	ErrOrderExistThis = errors.New("order whith this id of this user is exist")
	ErrOrderExistAnother = errors.New("order whith this id of another user is exist")
)
