package env

import (
	"strings"
)

// CompactMode controls which entries are removed during compaction.
type CompactMode string

const (
	// CompactNone disables compaction; all entries are passed through.
	CompactNone CompactMode = "none"
	// CompactDuplicates removes entries whose key appears more than once,
	// keeping only the last occurrence (same semantics as Dedupe).
	CompactDuplicates CompactMode = "duplicates"
	// CompactEmpty removes entries with an empty value.
	CompactEmpty CompactMode = "empty"
	// CompactBoth removes both duplicates and empty-value entries.
	CompactBoth CompactMode = "both"
)

// Compactor removes unwanted entries from an environment slice.
type Compactor struct {
	mode CompactMode
}

// NewCompactor returns a Compactor configured with the given mode.
// An empty mode string defaults to CompactNone.
func NewCompactor(mode CompactMode) *Compactor {
	if mode == "" {
		mode = CompactNone
	}
	return &Compactor{mode: mode}
}

// Apply filters the provided KEY=VALUE pairs according to the configured mode
// and returns the resulting slice. The input slice is never mutated.
func (c *Compactor) Apply(pairs []string) []string {
	if c.mode == CompactNone {
		return pairs
	}

	removeDupes := c.mode == CompactDuplicates || c.mode == CompactBoth
	removeEmpty := c.mode == CompactEmpty || c.mode == CompactBoth

	// First pass: collect last-seen index for each key (for dupe removal).
	lastSeen := make(map[string]int, len(pairs))
	if removeDupes {
		for i, p := range pairs {
			if idx := strings.IndexByte(p, '='); idx > 0 {
				lastSeen[p[:idx]] = i
			}
		}
	}

	out := make([]string, 0, len(pairs))
	for i, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			// Malformed entry — pass through unchanged.
			out = append(out, p)
			continue
		}

		key := p[:idx]
		value := p[idx+1:]

		if removeDupes && lastSeen[key] != i {
			continue
		}
		if removeEmpty && value == "" {
			continue
		}
		out = append(out, p)
	}

	if len(out) == 0 {
		return nil
	}
	return out
}
