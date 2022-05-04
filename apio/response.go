package apio

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/common-fate/apikit/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// JSON converts a Go value to JSON and sends it to the client.
// Under the hood, JSON uses logger.Get() to load a zap logger from context,
// so the logging middleware must run.
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

// Error sends an error reponse back to the client and logs the error internally
// if the error is of type io.Error we send it's message back to the client.
// Otherwise, we return a HTTP 500 code with an opaque response to avoid leaking any
// information from the server.
//
// Under the hood, Error uses logger.Get() to load a zap logger from context,
// so the logging middleware must run.
func Error(ctx context.Context, w http.ResponseWriter, err error) {
	// load the zap logger from context.
	log := logger.Get(ctx)

	log.Errorw("web handler error", zap.Error(err))

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
