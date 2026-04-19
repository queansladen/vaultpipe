package env

import (
	"strings"
)

// ScrubMode controls which entries are removed by the Scrubber.
type ScrubMode string

const (
	ScrubNone   ScrubMode = "none"
	ScrubEmpty  ScrubMode = "empty"
	ScrubPrefix ScrubMode = "prefix"
)

// Scrubber removes unwanted entries from an environment slice.
type Scrubber struct {
	mode     ScrubMode
	prefixes []string
}

// NewScrubber returns a Scrubber configured with the given mode and optional
// prefixes (used when mode is ScrubPrefix).
func NewScrubber(mode ScrubMode, prefixes ...string) *Scrubber {
	if mode == "" {
		mode = ScrubNone
	}
	return &Scrubber{mode: mode, prefixes: prefixes}
}

// Apply filters the provided env pairs according to the scrub mode.
func (s *Scrubber) Apply(pairs []string) []string {
	if s.mode == ScrubNone {
		return pairs
	}
	out := make([]string, 0, len(pairs))
	for _, p := range pairs {
		if s.shouldScrub(p) {
			continue
		}
		out = append(out, p)
	}
	return out
}

func (s *Scrubber) shouldScrub(pair string) bool {
	idx := strings.IndexByte(pair, '=')
	if idx < 0 {
		return false
	}
	switch s.mode {
	case ScrubEmpty:
		return pair[idx+1:] == ""
	case ScrubPrefix:
		key := pair[:idx]
		for _, pfx := range s.prefixes {
			if strings.HasPrefix(key, pfx) {
				return true
			}
		}
	}
	return false
}
