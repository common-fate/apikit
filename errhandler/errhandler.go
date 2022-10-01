package errhandler

import "context"

// Handler is an error handler which is called
// when apio.Error() is called in APIs.
type Handler interface {
	HandleError(err error)
}

var errHandlerKey = &contextKey{"errHandler"}

type contextKey struct {
	name string
}

// Set the error handler in context so that apio.Error()
// will log errors to it.
func Set(ctx context.Context, h Handler) context.Context {
	ctx = context.WithValue(ctx, errHandlerKey, h)
	return ctx
}

// Get retrieves an error handler from context.
// It returns nil if it hasn't been set.
func Get(ctx context.Context) Handler {
	if h, ok := ctx.Value(errHandlerKey).(Handler); ok {
		return h
	}
	return nil
}
