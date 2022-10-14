package apio

import (
	"context"
	"net/http"
)

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

// APIError is used to pass an error during the request through the
// application with web specific context.
type APIError struct {
	Err    error
	Status int
	Fields []FieldError
}

// NewRequestError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewRequestError(err error, status int) error {
	return &APIError{err, status, nil}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (e *APIError) Error() string {
	return e.Err.Error()
}

// RenderableErrors can render their own HTTP responses.
// They can be used to override the behaviour of apio.Error().
type RenderableError interface {
	RenderHTTP(ctx context.Context, w http.ResponseWriter)
}
