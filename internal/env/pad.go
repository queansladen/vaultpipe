package env

import (
	"fmt"
	"strings"
)

// PadMode controls how values are padded.
type PadMode string

const (
	PadModeNone  PadMode = "none"
	PadModeLeft  PadMode = "left"
	PadModeRight PadMode = "right"
)

// Padder pads environment variable values to a minimum width.
type Padder struct {
	mode  PadMode
	width int
	char  rune
}

// NewPadder creates a Padder that pads values to at least width characters
// using char as the fill rune. mode must be PadModeLeft or PadModeRight;
// PadModeNone or an unrecognised mode is a no-op.
func NewPadder(mode PadMode, width int, char rune) *Padder {
	if char == 0 {
		char = ' '
	}
	return &Padder{mode: mode, width: width, char: char}
}

// Apply pads the values of all well-formed KEY=VALUE pairs in pairs.
func (p *Padder) Apply(pairs []string) ([]string, error) {
	if p.mode == PadModeNone || p.mode == "" || p.width <= 0 {
		return pairs, nil
	}

	out := make([]string, len(pairs))
	for i, pair := range pairs {
		idx := strings.IndexByte(pair, '=')
		if idx < 0 {
			out[i] = pair
			continue
		}
		key := pair[:idx]
		val := pair[idx+1:]
		out[i] = fmt.Sprintf("%s=%s", key, p.pad(val))
	}
	return out, nil
}

func (p *Padder) pad(s string) string {
	delta := p.width - len(s)
	if delta <= 0 {
		return s
	}
	fill := strings.Repeat(string(p.char), delta)
	if p.mode == PadModeLeft {
		return fill + s
	}
	return s + fill
}
