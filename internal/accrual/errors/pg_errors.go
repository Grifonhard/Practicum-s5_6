package errors

import "errors"

var (
	ErrQueryExecution = errors.New("query execution failed")
	ErrConnectTimeout = errors.New("connect timeout")
)
