package env

import "fmt"

// LimitPolicy controls behaviour when the count exceeds Max.
type LimitPolicy int

const (
	LimitPolicyError LimitPolicy = iota
	LimitPolicyTruncate
)

// Limiter enforces a maximum number of environment variable pairs.
type Limiter struct {
	max    int
	policy LimitPolicy
}

// NewLimiter returns a Limiter. max <= 0 means no limit.
func NewLimiter(max int, policy LimitPolicy) *Limiter {
	return &Limiter{max: max, policy: policy}
}

// Apply enforces the limit on pairs. Returns an error when policy is
// LimitPolicyError and the count is exceeded.
func (l *Limiter) Apply(pairs []string) ([]string, error) {
	if l.max <= 0 || len(pairs) <= l.max {
		return pairs, nil
	}
	switch l.policy {
	case LimitPolicyTruncate:
		return pairs[:l.max], nil
	default:
		return nil, fmt.Errorf("env: pair count %d exceeds limit %d", len(pairs), l.max)
	}
}
