package env

import (
	"strings"
)

// Quoter serialises an env slice into shell-escaped KEY=VALUE lines.
type Quoter struct {
	shellSafe bool
}

// NewQuoter returns a Quoter. When shellSafe is true values containing
// special characters are wrapped in single quotes.
func NewQuoter(shellSafe bool) *Quoter {
	return &Quoter{shellSafe: shellSafe}
}

// Quote returns a copy of pairs with values optionally shell-quoted.
func (q *Quoter) Quote(pairs []string) []string {
	out := make([]string, 0, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out = append(out, p)
			continue
		}
		key := p[:idx]
		val := p[idx+1:]
		if q.shellSafe {
			val = shellQuote(val)
		}
		out = append(out, key+"="+val)
	}
	return out
}

// shellQuote wraps v in single quotes, escaping any embedded single quotes.
func shellQuote(v string) string {
	if !needsQuoting(v) {
		return v
	}
	return "'" + strings.ReplaceAll(v, "'", `'\''`) + "'"
}

func needsQuoting(v string) bool {
	for _, c := range v {
		if !isSafeChar(c) {
			return true
		}
	}
	return false
}

func isSafeChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_' || c == '-' || c == '.' || c == '/' || c == ':'
}
