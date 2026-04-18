package env

import (
	"os"
	"strings"
)

// Snapshot captures the current process environment as a key-value map.
type Snapshot struct {
	env map[string]string
}

// TakeSnapshot reads os.Environ and returns a Snapshot.
func TakeSnapshot() *Snapshot {
	return NewSnapshot(os.Environ())
}

// NewSnapshot builds a Snapshot from a slice of "KEY=VALUE" strings.
func NewSnapshot(pairs []string) *Snapshot {
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			continue
		}
		m[parts[0]] = parts[1]
	}
	return &Snapshot{env: m}
}

// Get returns the value for key and whether it was present.
func (s *Snapshot) Get(key string) (string, bool) {
	v, ok := s.env[key]
	return v, ok
}

// Keys returns all environment variable names in the snapshot.
func (s *Snapshot) Keys() []string {
	keys := make([]string, 0, len(s.env))
	for k := range s.env {
		keys = append(keys, k)
	}
	return keys
}

// Len returns the number of entries in the snapshot.
func (s *Snapshot) Len() int {
	return len(s.env)
}

// Pairs returns the snapshot as a slice of "KEY=VALUE" strings suitable
// for use with exec.Cmd.Env.
func (s *Snapshot) Pairs() []string {
	out := make([]string, 0, len(s.env))
	for k, v := range s.env {
		out = append(out, k+"="+v)
	}
	return out
}
