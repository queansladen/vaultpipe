package env

import (
	"fmt"
	"strings"
)

// FlattenMode controls how nested key segments are joined.
type FlattenMode string

const (
	FlattenNone       FlattenMode = "none"
	FlattenUnderscore FlattenMode = "underscore"
	FlattenDot        FlattenMode = "dot"
	FlattenDash       FlattenMode = "dash"
)

// Flattener joins multi-segment key names using a configurable separator,
// replacing the delimiter character found in keys (e.g. "/" or ".") with
// the chosen separator.
type Flattener struct {
	mode      FlattenMode
	separator string
	split     string
}

// NewFlattener returns a Flattener that replaces occurrences of splitOn in
// env-var keys with the separator implied by mode. If mode is FlattenNone
// or unrecognised, keys are passed through unchanged.
func NewFlattener(mode FlattenMode, splitOn string) *Flattener {
	sep := ""
	switch mode {
	case FlattenUnderscore:
		sep = "_"
	case FlattenDot:
		sep = "."
	case FlattenDash:
		sep = "-"
	}
	return &Flattener{mode: mode, separator: sep, split: splitOn}
}

// Apply returns a new slice with keys transformed according to the mode.
// Values are never modified. Malformed entries (no "=" separator) are
// passed through unchanged.
func (f *Flattener) Apply(pairs []string) ([]string, error) {
	if f.mode == FlattenNone || f.mode == "" || f.split == "" {
		out := make([]string, len(pairs))
		copy(out, pairs)
		return out, nil
	}

	out := make([]string, 0, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out = append(out, p)
			continue
		}
		key := p[:idx]
		val := p[idx+1:]
		newKey := strings.ReplaceAll(key, f.split, f.separator)
		out = append(out, fmt.Sprintf("%s=%s", newKey, val))
	}
	return out, nil
}
