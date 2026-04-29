package env

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// DecodeMode controls how base64 decoding is applied to env pairs.
type DecodeMode string

const (
	DecodeModeNone   DecodeMode = "none"
	DecodeModeValues DecodeMode = "values"
	DecodeModeKeys   DecodeMode = "keys"
	DecodeModeBoth   DecodeMode = "both"
)

// Decoder decodes base64-encoded keys and/or values in env pairs.
type Decoder struct {
	mode DecodeMode
}

// NewDecoder returns a Decoder for the given mode.
// An empty mode defaults to none (passthrough).
func NewDecoder(mode DecodeMode) *Decoder {
	if mode == "" {
		mode = DecodeModeNone
	}
	return &Decoder{mode: mode}
}

// Apply decodes base64 fields in each "KEY=VALUE" pair according to the mode.
// Pairs that fail decoding are passed through unchanged.
func (d *Decoder) Apply(pairs []string) ([]string, error) {
	if d.mode == DecodeModeNone {
		return pairs, nil
	}
	out := make([]string, 0, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out = append(out, p)
			continue
		}
		k := p[:idx]
		v := p[idx+1:]

		if d.mode == DecodeModeKeys || d.mode == DecodeModeBoth {
			if dec, err := base64.StdEncoding.DecodeString(k); err == nil {
				k = string(dec)
			}
		}
		if d.mode == DecodeModeValues || d.mode == DecodeModeBoth {
			if dec, err := base64.StdEncoding.DecodeString(v); err == nil {
				v = string(dec)
			}
		}
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}
	return out, nil
}
