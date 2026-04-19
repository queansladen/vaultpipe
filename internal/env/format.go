package env

import (
	"fmt"
	"strings"
)

// FormatMode controls how env pairs are serialised.
type FormatMode string

const (
	FormatModeRaw    FormatMode = "raw"    // KEY=VALUE
	FormatModeExport FormatMode = "export" // export KEY=VALUE
	FormatModeDotenv FormatMode = "dotenv" // KEY="VALUE"
)

// Formatter rewrites env pairs into a target format.
type Formatter struct {
	mode FormatMode
}

// NewFormatter returns a Formatter for the given mode.
// An empty or unrecognised mode defaults to FormatModeRaw.
func NewFormatter(mode FormatMode) *Formatter {
	switch mode {
	case FormatModeExport, FormatModeDotenv:
		return &Formatter{mode: mode}
	default:
		return &Formatter{mode: FormatModeRaw}
	}
}

// Apply rewrites each pair in pairs according to the formatter's mode.
// Malformed entries (no '=') are passed through unchanged.
func (f *Formatter) Apply(pairs []string) []string {
	out := make([]string, len(pairs))
	for i, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out[i] = p
			continue
		}
		key := p[:idx]
		val := p[idx+1:]
		switch f.mode {
		case FormatModeExport:
			out[i] = fmt.Sprintf("export %s=%s", key, val)
		case FormatModeDotenv:
			escaped := strings.ReplaceAll(val, `"`, `\"`)
			out[i] = fmt.Sprintf(`%s="%s"`, key, escaped)
		default:
			out[i] = p
		}
	}
	return out
}
