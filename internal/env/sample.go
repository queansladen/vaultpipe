package env

import (
	"math/rand"
	"strings"
)

// SampleMode controls how entries are selected.
type SampleMode string

const (
	SampleModeNone   SampleMode = "none"
	SampleModeRandom SampleMode = "random"
	SampleModeFirst  SampleMode = "first"
	SampleModeLast   SampleMode = "last"
)

// Sampler selects a subset of env entries.
type Sampler struct {
	mode SampleMode
	n    int
	rng  *rand.Rand
}

// NewSampler returns a Sampler that selects n entries using the given mode.
// A nil rng falls back to the default global source.
func NewSampler(mode SampleMode, n int, rng *rand.Rand) *Sampler {
	if mode == "" {
		mode = SampleModeNone
	}
	return &Sampler{mode: mode, n: n, rng: rng}
}

// Apply returns a sampled slice of env pairs.
func (s *Sampler) Apply(pairs []string) []string {
	if s.mode == SampleModeNone || s.n <= 0 || len(pairs) == 0 {
		return pairs
	}

	n := s.n
	if n > len(pairs) {
		n = len(pairs)
	}

	switch s.mode {
	case SampleModeFirst:
		out := make([]string, n)
		copy(out, pairs[:n])
		return out
	case SampleModeLast:
		out := make([]string, n)
		copy(out, pairs[len(pairs)-n:])
		return out
	case SampleModeRandom:
		copied := make([]string, len(pairs))
		copy(copied, pairs)
		if s.rng != nil {
			s.rng.Shuffle(len(copied), func(i, j int) { copied[i], copied[j] = copied[j], copied[i] })
		} else {
			rand.Shuffle(len(copied), func(i, j int) { copied[i], copied[j] = copied[j], copied[i] })
		}
		return copied[:n]
	default:
		return pairs
	}
}

// keyOf returns the key portion of a KEY=VALUE pair.
func sampleKeyOf(pair string) string {
	if idx := strings.IndexByte(pair, '='); idx >= 0 {
		return pair[:idx]
	}
	return pair
}

// Keys returns only the keys present in the sampled result.
func (s *Sampler) Keys(pairs []string) []string {
	sampled := s.Apply(pairs)
	keys := make([]string, 0, len(sampled))
	for _, p := range sampled {
		keys = append(keys, sampleKeyOf(p))
	}
	return keys
}
