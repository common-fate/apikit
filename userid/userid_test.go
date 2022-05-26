package userid

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	type testcase struct {
		name string
		uid  string
	}

	testcases := []testcase{
		{"ok", "testuser"},
		{"empty string", ""}, // should still get/set
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = Set(ctx, tc.uid)
			got := Get(ctx)
			assert.Equal(t, tc.uid, got)
		})
	}
}

func TestInit(t *testing.T) {
	type testcase struct {
		name string
		uid  string
	}

	testcases := []testcase{
		{"ok", "testuser"},
		{"empty string", ""}, // should still get/set
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = Init(ctx)
			Set(ctx, tc.uid)
			// we should be able to get the user ID from the original context after calling Init()
			got := Get(ctx)
			assert.Equal(t, tc.uid, got)
		})
	}
}
