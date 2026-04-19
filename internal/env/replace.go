package env

import "strings"

// Replacer performs find-and-replace operations on environment variable values.
type Replacer struct {
	rules []ReplaceRule
}

// ReplaceRule defines a single find-and-replace operation.
type ReplaceRule struct {
	Find    string
	Replace string
}

// NewReplacer creates a Replacer with the given rules.
func NewReplacer(rules []ReplaceRule) *Replacer {
	return &Replacer{rules: rules}
}

// Apply performs all replacements on each env pair in the slice.
// Pairs that do not contain a separator are passed through unchanged.
func (r *Replacer) Apply(pairs []string) []string {
	if len(r.rules) == 0 {
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
		for _, rule := range r.rules {
			if rule.Find == "" {
				continue
			}
			val = strings.ReplaceAll(val, rule.Find, rule.Replace)
		}
		out[i] = key + "=" + val
	}
	return out
}
