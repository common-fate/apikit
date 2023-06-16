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

type UnauthorisedError struct{}

// Error implements the error interface.
func (e UnauthorisedError) Error() string {
	return http.StatusText(http.StatusUnauthorized)
}

func Unauthorised() error {
	return UnauthorisedError{}
}

func IsUnauthorised(err error) bool {
	_, ok := err.(UnauthorisedError)
	return ok
}

type ForbiddenError struct{}

// Error implements the error interface.
func (e ForbiddenError) Error() string {
	return http.StatusText(http.StatusForbidden)
}

func Forbidden() error {
	return ForbiddenError{}
}

func IsForbidden(err error) bool {
	_, ok := err.(ForbiddenError)
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
