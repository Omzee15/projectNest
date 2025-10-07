package utils

import (
	"errors"
	"net/http"
)

// Custom error types
var (
	ErrNotFound     = errors.New("resource not found")
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrInternal     = errors.New("internal server error")
	ErrConflict     = errors.New("resource conflict")
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Err        error
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

// Error constructors
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Err:        ErrNotFound,
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Err:        ErrBadRequest,
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func NewInternalError(message string) *AppError {
	return &AppError{
		Err:        ErrInternal,
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Err:        ErrConflict,
		StatusCode: http.StatusConflict,
		Message:    message,
	}
}
