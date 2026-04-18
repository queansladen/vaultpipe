package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Handler listens for OS signals and cancels a context.
type Handler struct {
	signals []os.Signal
}

// NewHandler creates a Handler that watches the given signals.
// If no signals are provided, SIGINT and SIGTERM are used.
func NewHandler(sigs ...os.Signal) *Handler {
	if len(sigs) == 0 {
		sigs = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	return &Handler{signals: sigs}
}

// WithContext returns a derived context that is cancelled when one of the
// watched signals is received. The stop function should be called to release
// resources when the context is no longer needed.
func (h *Handler) WithContext(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, h.signals...)

	go func() {
		defer signal.Stop(ch)
		select {
		case <-ch:
			cancel()
		case <-parent.Done():
			cancel()
		}
	}()

	return ctx, cancel
}
