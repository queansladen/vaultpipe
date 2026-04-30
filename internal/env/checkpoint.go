package env

import "fmt"

// Checkpoint captures a named snapshot of an env slice for diff/debug purposes.
type Checkpoint struct {
	name    string
	records map[string]string
}

// NewCheckpoint records the current state of pairs under a given name.
func NewCheckpoint(name string, pairs []string) *Checkpoint {
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		k, v, ok := splitPair(p)
		if !ok {
			continue
		}
		m[k] = v
	}
	return &Checkpoint{name: name, records: m}
}

// Name returns the checkpoint label.
func (c *Checkpoint) Name() string { return c.name }

// Len returns the number of recorded keys.
func (c *Checkpoint) Len() int { return len(c.records) }

// Get returns the value for key and whether it was present in the checkpoint.
func (c *Checkpoint) Get(key string) (string, bool) {
	v, ok := c.records[key]
	return v, ok
}

// Diff returns keys whose values differ between c and next, as "KEY=old->new" strings.
func (c *Checkpoint) Diff(next *Checkpoint) []string {
	seen := make(map[string]struct{})
	var out []string

	for k, v := range next.records {
		seen[k] = struct{}{}
		if old, exists := c.records[k]; !exists {
			out = append(out, fmt.Sprintf("%s=<unset>->%s", k, v))
		} else if old != v {
			out = append(out, fmt.Sprintf("%s=%s->%s", k, old, v))
		}
	}

	for k, v := range c.records {
		if _, ok := seen[k]; !ok {
			out = append(out, fmt.Sprintf("%s=%s-><unset>", k, v))
		}
	}
	return out
}
