package env

import (
	"os"
	"strings"
)

// InheritMode controls how the parent environment is inherited.
type InheritMode int

const (
	// InheritAll passes the full parent environment to the child process.
	InheritAll InheritMode = iota
	// InheritNone starts with a completely clean environment.
	InheritNone
	// InheritFiltered passes only variables that match the supplied filter.
	InheritFiltered
)

// Inherit returns a slice of KEY=VALUE pairs derived from the current process
// environment according to the given mode and optional filter.
func Inherit(mode InheritMode, f *Filter) []string {
	switch mode {
	case InheritNone:
		return []string{}
	case InheritFiltered:
		if f == nil {
			return []string{}
		}
		return f.Apply(os.Environ())
	default: // InheritAll
		return append([]string(nil), os.Environ()...)
	}
}

// MergeEnv merges overlay pairs into base, with overlay taking precedence.
// Both slices are expected to contain KEY=VALUE strings.
func MergeEnv(base, overlay []string) []string {
	result := make(map[string]string, len(base)+len(overlay))
	order := make([]string, 0, len(base)+len(overlay))

	for _, pair := range base {
		k, v, _ := strings.Cut(pair, "=")
		if _, exists := result[k]; !exists {
			order = append(order, k)
		}
		result[k] = v
	}
	for _, pair := range overlay {
		k, v, _ := strings.Cut(pair, "=")
		if _, exists := result[k]; !exists {
			order = append(order, k)
		}
		result[k] = v
	}

	out := make([]string, 0, len(order))
	for _, k := range order {
		out = append(out, k+"="+result[k])
	}
	return out
}
