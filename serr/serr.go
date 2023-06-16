package serr

import (
	"errors"
	"net/http"
)

type SError struct {
	Err        error
	StatusCode int
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (e *SError) Error() string {
	return e.Err.Error()
}

func NotFound() error {
	return &SError{
		Err:        errors.New(http.StatusText(http.StatusNotFound)),
		StatusCode: http.StatusNotFound,
	}
}
