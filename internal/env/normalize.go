package env

import (
	"strings"
	"unicode"
)

// NormalizeMode controls how environment variable keys are normalized.
type NormalizeMode string

const (
	NormalizeModeNone  NormalizeMode = "none"
	NormalizeModeSnake NormalizeMode = "snake"
	NormalizeModeDot   NormalizeMode = "dot"
)

// Normalizer rewrites environment variable keys according to a normalization
// strategy, collapsing runs of non-alphanumeric characters into a separator.
type Normalizer struct {
	mode NormalizeMode
}

// NewNormalizer returns a Normalizer configured with the given mode.
// An empty or unrecognised mode defaults to NormalizeModeNone.
func NewNormalizer(mode NormalizeMode) *Normalizer {
	switch mode {
	case NormalizeModeSnake, NormalizeModeDot:
	default:
		mode = NormalizeModeNone
	}
	return &Normalizer{mode: mode}
}

// Apply normalizes the keys of each "KEY=VALUE" pair in pairs.
func (n *Normalizer) Apply(pairs []string) ([]string, error) {
	if n.mode == NormalizeModeNone {
		return pairs, nil
	}
	sep := "_"
	if n.mode == NormalizeModeDot {
		sep = "."
	}
	out := make([]string, 0, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out = append(out, p)
			continue
		}
		key := normalizeKey(p[:idx], sep)
		out = append(out, key+p[idx:])
	}
	return out, nil
}

// normalizeKey collapses runs of characters that are not letters or digits
// into sep, then trims leading/trailing separators.
func normalizeKey(key, sep string) string {
	var b strings.Builder
	inSep := false
	for _, r := range key {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			inSep = false
		} else {
			if !inSep && b.Len() > 0 {
				b.WriteString(sep)
			}
			inSep = true
		}
	}
	return strings.TrimSuffix(b.String(), sep)
}
