package serr

import (
	"net/http"
)

type NotFoundError struct{}

// Error implements the error interface.
func (e NotFoundError) Error() string {
	return http.StatusText(http.StatusNotFound)
}

func NotFound() error {
	return NotFoundError{}
}

func IsNotFound(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}

type BadRequestError struct {
	Msg string
}

// Error implements the error interface.
func (e BadRequestError) Error() string {
	return e.Msg
}

func BadRequest(msg string) error {
	return BadRequestError{Msg: msg}
}

func IsBadRequest(err error) bool {
	_, ok := err.(BadRequestError)
	return ok
}
