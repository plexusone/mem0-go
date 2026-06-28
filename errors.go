package mem0

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrNoAPIKey is returned when no API key is provided.
	ErrNoAPIKey = errors.New("mem0: API key is required")

	// ErrInvalidBackend is returned when an unsupported backend is specified.
	ErrInvalidBackend = errors.New("mem0: invalid backend")

	// ErrNotFound is returned when a resource is not found.
	ErrNotFound = errors.New("mem0: resource not found")

	// ErrUnauthorized is returned when authentication fails.
	ErrUnauthorized = errors.New("mem0: unauthorized")

	// ErrBadRequest is returned when the request is invalid.
	ErrBadRequest = errors.New("mem0: bad request")

	// ErrRateLimited is returned when the rate limit is exceeded.
	ErrRateLimited = errors.New("mem0: rate limit exceeded")
)

// APIError represents an error returned by the mem0 API.
type APIError struct {
	StatusCode int
	Detail     string
	Code       string
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("mem0: API error (status %d): %s", e.StatusCode, e.Detail)
	}
	return fmt.Sprintf("mem0: API error (status %d)", e.StatusCode)
}

// Unwrap returns the underlying error based on status code.
func (e *APIError) Unwrap() error {
	switch e.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusTooManyRequests:
		return ErrRateLimited
	default:
		return nil
	}
}

// IsNotFoundError returns true if the error indicates a resource was not found.
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsUnauthorizedError returns true if the error indicates an authentication failure.
func IsUnauthorizedError(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsBadRequestError returns true if the error indicates an invalid request.
func IsBadRequestError(err error) bool {
	return errors.Is(err, ErrBadRequest)
}

// IsRateLimitedError returns true if the error indicates rate limiting.
func IsRateLimitedError(err error) bool {
	return errors.Is(err, ErrRateLimited)
}
