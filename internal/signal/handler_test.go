package signal_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/signal"
)

func TestNewHandler_Defaults(t *testing.T) {
	h := signal.NewHandler()
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestWithContext_CancelledByParent(t *testing.T) {
	h := signal.NewHandler(syscall.SIGUSR1)

	parent, parentCancel := context.WithCancel(context.Background())
	ctx, stop := h.WithContext(parent)
	defer stop()

	parentCancel()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(500 * time.Millisecond):
		t.Fatal("context was not cancelled after parent cancel")
	}
}

func TestWithContext_CancelledBySignal(t *testing.T) {
	h := signal.NewHandler(syscall.SIGUSR1)

	ctx, stop := h.WithContext(context.Background())
	defer stop()

	// Send the signal to the current process.
	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1) //nolint:errcheck

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(500 * time.Millisecond):
		t.Fatal("context was not cancelled after signal")
	}
}

func TestWithContext_StopFuncReleasesResources(t *testing.T) {
	h := signal.NewHandler(syscall.SIGUSR2)
	_, stop := h.WithContext(context.Background())
	// Should not panic or block.
	stop()
}
