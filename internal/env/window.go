package env

import (
	"fmt"
	"strings"
)

// WindowMode controls which portion of the environment slice is returned.
type WindowMode string

const (
	WindowModeNone   WindowMode = "none"
	WindowModeHead   WindowMode = "head"
	WindowModeTail   WindowMode = "tail"
	WindowModeSlice  WindowMode = "slice"
)

// Windower returns a subset of environment pairs based on a configured window.
type Windower struct {
	mode   WindowMode
	offset int
	size   int
}

// NewWindower creates a Windower. offset and size are only used for WindowModeSlice.
// For head/tail, offset is ignored and size controls how many entries to keep.
func NewWindower(mode WindowMode, offset, size int) (*Windower, error) {
	switch mode {
	case WindowModeNone, WindowModeHead, WindowModeTail:
		// valid
	case WindowModeSlice:
		if offset < 0 {
			return nil, fmt.Errorf("window: offset must be >= 0, got %d", offset)
		}
		if size < 0 {
			return nil, fmt.Errorf("window: size must be >= 0, got %d", size)
		}
	default:
		mode = WindowModeNone
	}
	return &Windower{mode: mode, offset: offset, size: size}, nil
}

// Apply returns the windowed subset of pairs.
func (w *Windower) Apply(pairs []string) []string {
	if len(pairs) == 0 {
		return pairs
	}
	switch w.mode {
	case WindowModeHead:
		if w.size <= 0 || w.size >= len(pairs) {
			return pairs
		}
		return pairs[:w.size]
	case WindowModeTail:
		if w.size <= 0 || w.size >= len(pairs) {
			return pairs
		}
		return pairs[len(pairs)-w.size:]
	case WindowModeSlice:
		if w.offset >= len(pairs) {
			return nil
		}
		end := w.offset + w.size
		if w.size == 0 || end > len(pairs) {
			end = len(pairs)
		}
		return pairs[w.offset:end]
	default:
		return pairs
	}
}

// windowKey returns the key portion of a KEY=VALUE pair for display.
func windowKey(pair string) string {
	if idx := strings.IndexByte(pair, '='); idx >= 0 {
		return pair[:idx]
	}
	return pair
}

// ensure windowKey is referenced (used in tests via package).
var _ = windowKey
