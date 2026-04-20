package env

import (
	"fmt"
	"strings"
)

// ShellFormat controls the output format for shell evaluation.
type ShellFormat string

const (
	ShellFormatExport ShellFormat = "export"
	ShellFormatInline ShellFormat = "inline"
	ShellFormatUnset  ShellFormat = "unset"
)

// Sheller converts env pairs into shell-evaluable statements.
type Sheller struct {
	format ShellFormat
	quote  bool
}

// NewSheller returns a Sheller with the given format and quoting behaviour.
func NewSheller(format ShellFormat, quote bool) *Sheller {
	f := format
	if f == "" {
		f = ShellFormatExport
	}
	return &Sheller{format: f, quote: quote}
}

// Apply converts each "KEY=VALUE" pair into a shell statement.
// Malformed entries are passed through unchanged.
func (s *Sheller) Apply(pairs []string) []string {
	out := make([]string, 0, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out = append(out, p)
			continue
		}
		key := p[:idx]
		val := p[idx+1:]
		if s.quote {
			val = shellQuoteValue(val)
		}
		switch s.format {
		case ShellFormatInline:
			out = append(out, fmt.Sprintf("%s=%s", key, val))
		case ShellFormatUnset:
			out = append(out, fmt.Sprintf("unset %s", key))
		default:
			out = append(out, fmt.Sprintf("export %s=%s", key, val))
		}
	}
	return out
}

// shellQuoteValue wraps val in single quotes, escaping embedded single quotes.
func shellQuoteValue(val string) string {
	return "'" + strings.ReplaceAll(val, "'", "'\\''" ) + "'"
}
