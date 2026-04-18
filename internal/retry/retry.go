package retry

import (
	"context"
	"errors"
	"time"
)

// Policy defines retry behaviour.
type Policy struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// ErrMaxAttempts is returned when all attempts are exhausted.
var ErrMaxAttempts = errors.New("retry: max attempts reached")

// Do executes fn up to policy.MaxAttempts times, backing off between attempts.
// It stops early if ctx is cancelled or fn returns nil.
func Do(ctx context.Context, p Policy, fn func() error) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	delay := p.Delay
	var lastErr error
	for i := 0; i < p.MaxAttempts; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if i < p.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			if p.Multiplier > 0 {
				delay = time.Duration(float64(delay) * p.Multiplier)
			}
		}
	}
	return errors.Join(ErrMaxAttempts, lastErr)
}
