// Package errhandler defines a Handler interface for handling errors.
// When developing an API you can implement this interface if you'd like
// to use an error tracking service like Sentry.
//
// To use an error handler in your API, call errhandler.Set() as part of your
// middleware stack. You'll need to write a struct which implements the errhandler.Handler
// interface. Your HandleError method on the struct should contain all integration-specific
// logic to deal with the error, such as dispatching it to Sentry.
//
// When calling apio.Error(), if a Handler exists in the provided context, HandleError() will
// be called.
package errhandler
