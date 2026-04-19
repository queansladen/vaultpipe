package env

import "strings"

// MaskMode controls which part of an env value is masked.
type MaskMode string

const (
	MaskNone   MaskMode = "none"
	MaskFull   MaskMode = "full"
	MaskPartial MaskMode = "partial"
)

// Masker replaces sensitive values in env pairs.
type Masker struct {
	mode      MaskMode
	char      string
	keepChars int
}

// NewMasker returns a Masker with the given mode, mask character, and
// number of trailing characters to keep when mode is partial.
func NewMasker(mode MaskMode, char string, keepChars int) *Masker {
	if char == "" {
		char = "*"
	}
	if keepChars < 0 {
		keepChars = 0
	}
	return &Masker{mode: mode, char: char, keepChars: keepChars}
}

// Apply masks the value portion of each KEY=VALUE pair.
func (m *Masker) Apply(pairs []string) []string {
	if m.mode == MaskNone {
		return pairs
	}
	out := make([]string, len(pairs))
	for i, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out[i] = p
			continue
		}
		key := p[:idx]
		val := p[idx+1:]
		out[i] = key + "=" + m.maskValue(val)
	}
	return out
}

func (m *Masker) maskValue(val string) string {
	switch m.mode {
	case MaskFull:
		return strings.Repeat(m.char, len(val))
	case MaskPartial:
		if len(val) <= m.keepChars {
			return strings.Repeat(m.char, len(val))
		}
		masked := len(val) - m.keepChars
		return strings.Repeat(m.char, masked) + val[masked:]
	default:
		return val
	}
}
