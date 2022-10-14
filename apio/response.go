package apio

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/common-fate/apikit/errhandler"
	"github.com/common-fate/apikit/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// JSON converts a Go value to JSON and sends it to the client.
// Under the hood, JSON uses logger.Get() to load a zap logger from the provided context.
func JSON(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) {
	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return
	}

	// load the zap logger from context.
	log := logger.Get(ctx)

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Errorw("marshalling JSON", zap.Error(err))
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		log.Errorw("writing response", zap.Error(err))
	}
}

// Error sends an error reponse back to the client and logs the error internally.
// If the error is of type apio.Error we will send the error message back to the client.
// Otherwise, we return a HTTP 500 code with an opaque response to avoid leaking any
// information from the server.
//
// Under the hood, Error uses logger.Get() to load a zap logger from the provided context.
//
// The response body is in the format:
//
//	{"error": "msg"}
//
// If errhandler.Handler is set in the context, it will always be called with the error.
// You can check the error type in your error handler to determine the status code of the error.
func Error(ctx context.Context, w http.ResponseWriter, err error) {
	// load the zap logger from context.
	log := logger.Get(ctx)

	// dispatch an error if we have an error handler we can send it to.
	if h := errhandler.Get(ctx); h != nil {
		h.HandleError(err)
	}

	log.Errorw("web handler error", zap.Error(err))

	// If the error implements RenderableError, it can render a custom response.
	if re, ok := err.(RenderableError); ok {
		re.RenderHTTP(ctx, w)
		return
	}

	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := errors.Cause(err).(*APIError); ok {
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		JSON(ctx, w, er, webErr.Status)
		return
	}

	// If not, the handler sent any arbitrary error value so use 500.
	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	JSON(ctx, w, er, http.StatusInternalServerError)
}

// ErrorString sends an error response designated status code and error message.
// The response body is in the format:
//
//	{"error": "msg"}
//
// It's a convenience wrapper over apio.Error().
func ErrorString(ctx context.Context, w http.ResponseWriter, msg string, code int) {
	err := NewRequestError(errors.New(msg), code)
	Error(ctx, w, err)
}
