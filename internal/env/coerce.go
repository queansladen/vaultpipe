package env

import (
	"fmt"
	"strings"
)

// CoerceMode controls how values are coerced before injection.
type CoerceMode int

const (
	CoerceNone  CoerceMode = iota
	CoerceLower            // lowercase the value
	CoerceUpper            // uppercase the value
)

// Coercer transforms env pair values according to a CoerceMode.
type Coercer struct {
	mode CoerceMode
}

// NewCoercer returns a Coercer for the given mode.
func NewCoercer(mode CoerceMode) *Coercer {
	return &Coercer{mode: mode}
}

// Coerce applies the coercion mode to a slice of KEY=VALUE pairs.
// Malformed entries (no "=") are passed through unchanged.
func (c *Coercer) Coerce(pairs []string) []string {
	if c.mode == CoerceNone {
		out := make([]string, len(pairs))
		copy(out, pairs)
		return out
	}
	out := make([]string, 0, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out = append(out, p)
			continue
		}
		key := p[:idx]
		val := p[idx+1:]
		switch c.mode {
		case CoerceLower:
			val = strings.ToLower(val)
		case CoerceUpper:
			val = strings.ToUpper(val)
		}
		out = append(out, fmt.Sprintf("%s=%s", key, val))
	}
	return out
}
