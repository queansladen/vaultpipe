package env

import (
	"os"
	"strings"
)

// Expander resolves ${VAR} and $VAR references within env values
// using a provided lookup map, falling back to os.Getenv.
type Expander struct {
	lookup map[string]string
	fallback bool
}

// NewExpander creates an Expander backed by the given map.
// If fallback is true, unknown variables are resolved via os.Getenv.
func NewExpander(lookup map[string]string, fallback bool) *Expander {
	return &Expander{lookup: lookup, fallback: fallback}
}

// Expand replaces variable references in s.
func (e *Expander) Expand(s string) string {
	return os.Expand(s, e.resolve)
}

// ExpandAll returns a new slice with every value expanded.
func (e *Expander) ExpandAll(pairs []string) []string {
	out := make([]string, len(pairs))
	for i, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out[i] = p
			continue
		}
		out[i] = p[:idx+1] + e.Expand(p[idx+1:])
	}
	return out
}

func (e *Expander) resolve(key string) string {
	if v, ok := e.lookup[key]; ok {
		return v
	}
	if e.fallback {
		return os.Getenv(key)
	}
	return ""
}
