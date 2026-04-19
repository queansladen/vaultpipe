package env

import (
	"encoding/base64"
	"strings"
)

// EncodeMode controls how env values are encoded.
type EncodeMode string

const (
	EncodeModeNone   EncodeMode = "none"
	EncodeModeBase64 EncodeMode = "base64"
)

// Encoder applies an encoding transformation to env var values.
type Encoder struct {
	mode EncodeMode
}

// NewEncoder returns an Encoder for the given mode.
// An empty or unrecognised mode defaults to EncodeModeNone.
func NewEncoder(mode EncodeMode) *Encoder {
	switch mode {
	case EncodeModeBase64:
		return &Encoder{mode: EncodeModeBase64}
	default:
		return &Encoder{mode: EncodeModeNone}
	}
}

// Apply encodes the value portion of each "KEY=VALUE" pair.
func (e *Encoder) Apply(pairs []string) []string {
	if e.mode == EncodeModeNone {
		out := make([]string, len(pairs))
		copy(out, pairs)
		return out
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
		encoded := base64.StdEncoding.EncodeToString([]byte(val))
		out = append(out, key+"="+encoded)
	}
	return out
}
