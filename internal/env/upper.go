package env

import (
	"strings"
)

// CaseMode controls how environment variable keys or values are cased.
type CaseMode string

const (
	CaseModeNone  CaseMode = "none"
	CaseModeUpper CaseMode = "upper"
	CaseModeLower CaseMode = "lower"
)

// NewCaser returns a pipeline step that applies the given CaseMode to
// environment variable keys, values, or both.
func NewCaser(target, mode CaseMode) func([]string) ([]string, error) {
	return func(pairs []string) ([]string, error) {
		if mode == CaseModeNone || mode == "" {
			return pairs, nil
		}
		apply := strings.ToUpper
		if mode == CaseModeLower {
			apply = strings.ToLower
		}
		out := make([]string, len(pairs))
		for i, p := range pairs {
			idx := strings.IndexByte(p, '=')
			if idx < 0 {
				out[i] = p
				continue
			}
			k := p[:idx]
			v := p[idx+1:]
			switch target {
			case "key":
				k = apply(k)
			case "value":
				v = apply(v)
			default:
				k = apply(k)
				v = apply(v)
			}
			out[i] = k + "=" + v
		}
		return out, nil
	}
}
