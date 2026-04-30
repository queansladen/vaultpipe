package env

import (
	"fmt"
	"os"
	"strings"
)

// Interpolator expands ${VAR} and $VAR references within env-pair values
// using a precedence chain: overlay map → current pairs → OS environment.
type Interpolator struct {
	overlay map[string]string
	fallback bool
}

// NewInterpolator returns an Interpolator that resolves variable references
// inside values. When fallback is true, unresolved references fall through to
// os.Getenv; otherwise they expand to an empty string.
func NewInterpolator(overlay map[string]string, fallback bool) *Interpolator {
	if overlay == nil {
		overlay = make(map[string]string)
	}
	return &Interpolator{overlay: overlay, fallback: fallback}
}

// Apply expands variable references in the value portion of every pair.
func (i *Interpolator) Apply(pairs []string) ([]string, error) {
	// Build a lookup from the pairs themselves so earlier entries can be
	// referenced by later ones.
	local := make(map[string]string, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			continue
		}
		local[p[:idx]] = p[idx+1:]
	}

	lookup := func(key string) string {
		if v, ok := i.overlay[key]; ok {
			return v
		}
		if v, ok := local[key]; ok {
			return v
		}
		if i.fallback {
			return os.Getenv(key)
		}
		return ""
	}

	out := make([]string, 0, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out = append(out, p)
			continue
		}
		key := p[:idx]
		val := os.Expand(p[idx+1:], lookup)
		out = append(out, fmt.Sprintf("%s=%s", key, val))
	}
	return out, nil
}
