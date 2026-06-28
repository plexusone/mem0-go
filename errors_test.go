package mem0

import (
	"errors"
	"net/http"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		expected string
	}{
		{
			name:     "with detail",
			err:      &APIError{StatusCode: 400, Detail: "invalid request"},
			expected: "mem0: API error (status 400): invalid request",
		},
		{
			name:     "without detail",
			err:      &APIError{StatusCode: 500},
			expected: "mem0: API error (status 500)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestAPIError_Unwrap(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		expected error
	}{
		{
			name:     "not found",
			err:      &APIError{StatusCode: http.StatusNotFound},
			expected: ErrNotFound,
		},
		{
			name:     "unauthorized",
			err:      &APIError{StatusCode: http.StatusUnauthorized},
			expected: ErrUnauthorized,
		},
		{
			name:     "bad request",
			err:      &APIError{StatusCode: http.StatusBadRequest},
			expected: ErrBadRequest,
		},
		{
			name:     "rate limited",
			err:      &APIError{StatusCode: http.StatusTooManyRequests},
			expected: ErrRateLimited,
		},
		{
			name:     "unknown",
			err:      &APIError{StatusCode: 500},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unwrapped := tt.err.Unwrap()
			if unwrapped != tt.expected {
				t.Errorf("got %v, want %v", unwrapped, tt.expected)
			}
		})
	}
}

func TestErrorHelpers(t *testing.T) {
	notFoundErr := &APIError{StatusCode: http.StatusNotFound}
	unauthorizedErr := &APIError{StatusCode: http.StatusUnauthorized}
	badRequestErr := &APIError{StatusCode: http.StatusBadRequest}
	rateLimitedErr := &APIError{StatusCode: http.StatusTooManyRequests}

	if !IsNotFoundError(notFoundErr) {
		t.Error("expected IsNotFoundError to return true")
	}
	if !errors.Is(notFoundErr, ErrNotFound) {
		t.Error("expected errors.Is to match ErrNotFound")
	}

	if !IsUnauthorizedError(unauthorizedErr) {
		t.Error("expected IsUnauthorizedError to return true")
	}

	if !IsBadRequestError(badRequestErr) {
		t.Error("expected IsBadRequestError to return true")
	}

	if !IsRateLimitedError(rateLimitedErr) {
		t.Error("expected IsRateLimitedError to return true")
	}
}
