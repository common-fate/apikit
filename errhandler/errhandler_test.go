package errhandler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testHandler struct{}

func (h *testHandler) HandleError(err error) {}

func TestContext(t *testing.T) {
	h := &testHandler{}
	ctx := context.Background()
	ctx = Set(ctx, h)
	got := Get(ctx)
	assert.Equal(t, h, got)
}
