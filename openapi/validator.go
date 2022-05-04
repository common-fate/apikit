// Package openapi contains middleware to validate requests against an OpenAPI schema.
package openapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/common-fate/apikit/apio"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

// Options to customize request validation, openapi3filter specified options will be passed through.
type Options struct {
	Options openapi3filter.Options
}

// Validator creates middleware to validate request by swagger spec.
// This middleware is good for net/http either since go-chi is 100% compatible with net/http.
//
// This code is taken from the oapi-codegen package and has been modified to return JSON error
// responses rather than plaintext ones.
func Validator(swagger *openapi3.T) func(next http.Handler) http.Handler {
	return ValidatorWithOptions(swagger, nil)
}

// ValidatorWithOptions Creates middleware to validate request by swagger spec.
// This middleware is good for net/http either since go-chi is 100% compatible with net/http.
//
// This code is taken from the oapi-codegen package and has been modified to return JSON error
// responses rather than plaintext ones.
func ValidatorWithOptions(swagger *openapi3.T, options *Options) func(next http.Handler) http.Handler {
	router, err := gorillamux.NewRouter(swagger)
	if err != nil {
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// validate request
			if err := validateRequest(r, router, options); err != nil {
				apio.Error(r.Context(), w, err)
				return
			}

			// serve
			next.ServeHTTP(w, r)
		})
	}

}

// This function is called from the middleware above and actually does the work
// of validating a request.
func validateRequest(r *http.Request, router routers.Router, options *Options) error {

	// Find route
	route, pathParams, err := router.FindRoute(r)
	if err != nil {
		// We failed to find a matching route for the request.
		return &apio.APIError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
	}

	if options != nil {
		requestValidationInput.Options = &options.Options
	}

	if err := openapi3filter.ValidateRequest(context.Background(), requestValidationInput); err != nil {
		switch e := err.(type) {
		case *openapi3filter.RequestError:
			// We've got a bad request
			// Split up the verbose error by lines and return the first one
			// openapi errors seem to be multi-line with a decent message on the first
			errorLines := strings.Split(e.Error(), "\n")
			return &apio.APIError{
				Err:    fmt.Errorf(errorLines[0]),
				Status: http.StatusBadRequest,
			}
		case *openapi3filter.SecurityRequirementsError:
			return &apio.APIError{
				Err:    err,
				Status: http.StatusUnauthorized,
			}
		default:
			// This should never happen today, but if our upstream code changes,
			// we don't want to crash the server, so handle the unexpected error.
			return &apio.APIError{
				Err:    fmt.Errorf("error validating route: %s", err.Error()),
				Status: http.StatusInternalServerError,
			}
		}
	}

	return nil
}
