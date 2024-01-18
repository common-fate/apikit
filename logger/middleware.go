package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/common-fate/apikit/userid"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logCtxKey = &contextKey{"log"}

type contextKey struct {
	name string
}

// Middleware is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func Middleware(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()

			ctx := r.Context()
			reqID := middleware.GetReqID(ctx)

			// add the logger to context, so that logger.Get() can be used to retrieve it in
			// API endpoints.
			ctx = context.WithValue(ctx, logCtxKey, l.With(zap.String("reqId", reqID)).Sugar())

			// init the user ID on the context.
			// children middleware further down the stack can write the user ID to it.
			ctx = userid.Init(ctx)

			r = r.WithContext(ctx)

			defer func() {
				fields := []zapcore.Field{
					zap.String("proto", r.Proto),
					zap.String("remote", r.RemoteAddr),
					zap.String("request", r.RequestURI),
					zap.String("method", r.Method),
					zap.Duration("took", time.Since(t1)),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
					zap.String("reqId", reqID),
				}

				// get the user ID from the request context.
				// Authentication middleware may run *after* our logging
				// middleware, so we call it after next.ServeHTTP is complete.
				uid := userid.Get(ctx)
				if uid != "" {
					fields = append(fields, zap.String("userId", uid))
				}

				l.Info("Served", fields...)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// Get returns the logger in context, if there is one.
// If there isn't, it returns the global logger.
func Get(ctx context.Context) *zap.SugaredLogger {
	if l, ok := ctx.Value(logCtxKey).(*zap.SugaredLogger); ok {
		return l
	}

	return zap.S()
}

// Set the logger in the context.
func Set(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, logCtxKey, logger)
}
