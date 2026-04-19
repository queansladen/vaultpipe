package env

import (
	"fmt"
	"strings"
)

// Joiner merges two env slices using a configurable separator for duplicate keys.
type Joiner struct {
	separator string
}

// NewJoiner returns a Joiner that concatenates values for duplicate keys with sep.
func NewJoiner(sep string) *Joiner {
	if sep == "" {
		sep = ","
	}
	return &Joiner{separator: sep}
}

// Join merges base and overlay slices. When the same key appears in both,
// the values are concatenated with the configured separator rather than
// one overwriting the other.
func (j *Joiner) Join(base, overlay []string) []string {
	baseMap := make(map[string]string, len(base))
	order := make([]string, 0, len(base))

	for _, pair := range base {
		idx := strings.IndexByte(pair, '=')
		if idx < 0 {
			continue
		}
		k, v := pair[:idx], pair[idx+1:]
		if _, exists := baseMap[k]; !exists {
			order = append(order, k)
		}
		baseMap[k] = v
	}

	for _, pair := range overlay {
		idx := strings.IndexByte(pair, '=')
		if idx < 0 {
			continue
		}
		k, v := pair[:idx], pair[idx+1:]
		if existing, exists := baseMap[k]; exists {
			baseMap[k] = existing + j.separator + v
		} else {
			order = append(order, k)
			baseMap[k] = v
		}
	}

	out := make([]string, 0, len(order))
	for _, k := range order {
		out = append(out, fmt.Sprintf("%s=%s", k, baseMap[k]))
	}
	return out
}
