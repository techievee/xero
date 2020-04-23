package xeroErrors

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewUnexpectedGenericError(t *testing.T) {
	tests := []struct {
		name string
		args error
		code int
	}{
		{
			"error is a normal error message, output should be UnexpectedError",
			errors.New("database error"),
			http.StatusInternalServerError,
		},
		{
			"error is a forbidden error, output should be ForbiddenError",
			XeroForbiddenError(),
			http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUnexpectedGenericError(tt.args); tt.code != got.Code {
				t.Errorf("NewUnexpectedGenericError() = %v, want %v", got.Code, tt.code)
			}
		})
	}
}
