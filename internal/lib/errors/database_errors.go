package errors

import "errors"

const (
	ErrPostgresUniqueViolation = "23505"
)

var (
	ErrConnectTimeout = errors.New("connect timeout")
)
