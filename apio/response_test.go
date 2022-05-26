package apio

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorString(t *testing.T) {
	type testcase struct {
		name string
		msg  string
		code int
		want string
	}

	testcases := []testcase{
		{name: "ok", msg: "test", code: http.StatusBadRequest, want: `{"error":"test"}`},
		{name: "empty msg", msg: "", code: http.StatusBadRequest, want: `{"error":""}`},
	}

	ctx := context.Background()

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()

			ErrorString(ctx, rr, tc.msg, tc.code)

			assert.Equal(t, tc.code, rr.Code)

			data, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.want, string(data))
		})
	}
}
