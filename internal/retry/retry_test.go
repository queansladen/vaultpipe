package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/retry"
)

var errTemp = errors.New("temporary error")

func fastPolicy(attempts int) retry.Policy {
	return retry.Policy{MaxAttempts: attempts, Delay: time.Millisecond, Multiplier: 1.0}
}

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastPolicy(3), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesUntilSuccess(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastPolicy(3), func() error {
		calls++
		if calls < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastPolicy(3), func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, retry.ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := retry.Do(ctx, fastPolicy(5), func() error {
		return errTemp
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDefaultPolicy(t *testing.T) {
	p := retry.DefaultPolicy()
	if p.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts 3, got %d", p.MaxAttempts)
	}
	if p.Multiplier != 2.0 {
		t.Errorf("expected Multiplier 2.0, got %f", p.Multiplier)
	}
}
