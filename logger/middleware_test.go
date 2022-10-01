package logger

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/common-fate/apikit/userid"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestMiddlewareUserID(t *testing.T) {
	// source: https://stackoverflow.com/questions/70400426/how-to-properly-capture-zap-logger-output-in-unit-tests
	mycore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stderr,
		zapcore.InfoLevel,
	)
	// test core
	observed, logs := observer.New(zapcore.InfoLevel)

	// new logger with the two cores tee'd together
	logger := zap.New(zapcore.NewTee(mycore, observed))

	m := Middleware(logger)

	r := chi.NewRouter()
	r.Use(testRequestInfo)
	r.Use(m)
	r.Use(testUserID)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	want := []zapcore.Field{
		zap.String("reqId", "req-123"),
		zap.String("userId", "usr-123"),
	}
	for _, w := range want {
		assert.Contains(t, logs.All()[0].Context, w)
	}
}

func testRequestInfo(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, middleware.RequestIDKey, "req-123")
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func testUserID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = userid.Set(ctx, "usr-123")
		r = r.Clone(ctx)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
