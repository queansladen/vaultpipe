package lease

import (
	"context"
	"log"
	"time"
)

// SecretLease holds metadata about a Vault dynamic secret lease.
type SecretLease struct {
	LeaseID   string
	Duration  time.Duration
	Renewable bool
}

// Renewer periodically renews a Vault lease before it expires.
type Renewer struct {
	renewFn  func(ctx context.Context, leaseID string) (time.Duration, error)
	lease    SecretLease
	threshold float64 // fraction of lease duration at which to renew
}

// NewRenewer creates a Renewer for the given lease.
// threshold is a value between 0 and 1 (e.g. 0.75 renews at 75% of lease duration).
func NewRenewer(lease SecretLease, threshold float64, renewFn func(ctx context.Context, leaseID string) (time.Duration, error)) *Renewer {
	if threshold <= 0 || threshold >= 1 {
		threshold = 0.75
	}
	return &Renewer{
		renewFn:   renewFn,
		lease:     lease,
		threshold: threshold,
	}
}

// Start begins the renewal loop, blocking until ctx is cancelled or the lease is non-renewable.
func (r *Renewer) Start(ctx context.Context) {
	if !r.lease.Renewable || r.lease.LeaseID == "" {
		return
	}

	for {
		waitDuration := time.Duration(float64(r.lease.Duration) * r.threshold)
		select {
		case <-ctx.Done():
			return
		case <-time.After(waitDuration):
			newDuration, err := r.renewFn(ctx, r.lease.LeaseID)
			if err != nil {
				log.Printf("vaultpipe: lease renewal failed for %s: %v", r.lease.LeaseID, err)
				return
			}
			r.lease.Duration = newDuration
		}
	}
}
