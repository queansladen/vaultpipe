package vault

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultpipe/internal/retry"
)

// isRetryable returns true for transient HTTP status codes.
func isRetryable(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}
	return false
}

// ReadSecretDataWithRetry wraps ReadSecretData with retry logic.
func (c *Client) ReadSecretDataWithRetry(ctx context.Context, path string, p retry.Policy) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := retry.Do(ctx, p, func() error {
		var innerErr error
		data, innerErr = c.ReadSecretData(ctx, path)
		if innerErr != nil {
			if re, ok := innerErr.(*RetryableError); ok && isRetryable(re.StatusCode) {
				return innerErr
			}
			// Non-retryable: wrap to signal stop.
			return &fatalError{cause: innerErr}
		}
		return nil
	})
	if err != nil {
		if fe, ok := err.(*fatalError); ok {
			return nil, fe.cause
		}
		return nil, err
	}
	return data, nil
}

// RetryableError carries an HTTP status code for retry decisions.
type RetryableError struct {
	StatusCode int
	Msg        string
}

func (e *RetryableError) Error() string {
	return fmt.Sprintf("vault: retryable error %d: %s", e.StatusCode, e.Msg)
}

type fatalError struct{ cause error }

func (e *fatalError) Error() string { return e.cause.Error() }
