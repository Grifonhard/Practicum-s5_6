package accrual

import (
	"errors"
)

var (
	ErrStatusNoContent           = errors.New("order not registered in the accrual system")
	ErrStatusInternalServerError = errors.New("internal server error")
	ErrTooManyRequests           = errors.New("too many requests")
	ErrUnexpected                = errors.New("unexpected response status")
)
