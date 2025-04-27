package evmrpc

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWSConnectionHandler(t *testing.T) {
	type args struct {
		handler http.Handler
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "wraps handler",
			args: args{
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWSConnectionHandler(tt.args.handler)
			// Check the return type
			wsHandler, ok := got.(*wsConnectionHandler)
			assert.True(t, ok, "should return *wsConnectionHandler")
			// Check if the underlying handler is the same function
			assert.Equal(t,
				reflect.ValueOf(tt.args.handler).Pointer(),
				reflect.ValueOf(wsHandler.underlying).Pointer(),
				"underlying handler should be the same function",
			)
		})
	}
}

func Test_wsConnectionHandler_ServeHTTP(t *testing.T) {
	called := false
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusTeapot)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	rr := &mockResponseWriter{header: http.Header{}}

	h := &wsConnectionHandler{
		underlying: mockHandler,
	}
	h.ServeHTTP(rr, req)

	assert.True(t, called, "underlying handler should be called")
	assert.Equal(t, http.StatusTeapot, rr.status)
}

// mockResponseWriter is used to capture the status code written
type mockResponseWriter struct {
	header http.Header
	status int
}

func (m *mockResponseWriter) Header() http.Header {
	return m.header
}
func (m *mockResponseWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.status = statusCode
}
