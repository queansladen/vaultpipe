package lease

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewRenewer_DefaultsThreshold(t *testing.T) {
	lease := SecretLease{LeaseID: "id", Duration: time.Minute, Renewable: true}
	r := NewRenewer(lease, 0, nil)
	if r.threshold != 0.75 {
		t.Errorf("expected default threshold 0.75, got %v", r.threshold)
	}
}

func TestStart_NonRenewable_Returns(t *testing.T) {
	lease := SecretLease{LeaseID: "id", Duration: time.Minute, Renewable: false}
	r := NewRenewer(lease, 0.75, func(ctx context.Context, id string) (time.Duration, error) {
		t.Fatal("renewFn should not be called for non-renewable lease")
		return 0, nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	done := make(chan struct{})
	go func() {
		r.Start(ctx)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Start did not return for non-renewable lease")
	}
}

func TestStart_CancelledContext(t *testing.T) {
	lease := SecretLease{LeaseID: "abc", Duration: 10 * time.Second, Renewable: true}
	r := NewRenewer(lease, 0.75, func(ctx context.Context, id string) (time.Duration, error) {
		return 10 * time.Second, nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		r.Start(ctx)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Start did not return after context cancellation")
	}
}

func TestStart_CallsRenewFn(t *testing.T) {
	var calls atomic.Int32
	lease := SecretLease{LeaseID: "xyz", Duration: 100 * time.Millisecond, Renewable: true}
	ctx, cancel := context.WithCancel(context.Background())
	r := NewRenewer(lease, 0.5, func(ctx context.Context, id string) (time.Duration, error) {
		calls.Add(1)
		cancel()
		return 100 * time.Millisecond, nil
	})
	r.Start(ctx)
	if calls.Load() < 1 {
		t.Error("expected renewFn to be called at least once")
	}
}
