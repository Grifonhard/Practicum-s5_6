package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type Type string

const (
	BadRequest         Type = "BAD_REQUEST"         // Validation errors / BadInput
	Internal           Type = "INTERNAL"            // Server (500) and fallback errors
	NotFound           Type = "NOT_FOUND"           // For not finding resource
	ServiceUnavailable Type = "SERVICE_UNAVAILABLE" // For long-running handlers
	Conflict           Type = "CONFLICT"
)

type Error struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Status() int {
	switch e.Type {
	case BadRequest:
		return http.StatusBadRequest
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

// NewBadRequest to create 400 errors (validation, for example)
func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: fmt.Sprintf("Bad request. Reason: %v", reason),
	}
}

// NewInternal for 500 errors and unknown errors
func NewInternal() *Error {
	return &Error{
		Type:    Internal,
		Message: "Internal server error.",
	}
}

// NewNotFound to create an error for 404
func NewNotFound(name string, value string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("resource: %v with value: %v not found", name, value),
	}
}

// NewServiceUnavailable to create an error for 503
func NewServiceUnavailable() *Error {
	return &Error{
		Type:    ServiceUnavailable,
		Message: "Service unavailable or timed out",
	}
}

// NewConflict to create an error for 409
func NewConflict(err error) *Error {
	return &Error{
		Type:    Conflict,
		Message: err.Error(),
	}
}
