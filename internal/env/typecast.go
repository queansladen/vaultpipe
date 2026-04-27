package env

import (
	"fmt"
	"strconv"
	"strings"
)

// CastMode controls how environment variable values are type-coerced before
// being written into the process environment.
type CastMode string

const (
	CastNone    CastMode = "none"
	CastBool    CastMode = "bool"
	CastNumeric CastMode = "numeric"
	CastAuto    CastMode = "auto"
)

// Typecaster normalises values to a canonical string representation of their
// inferred or declared type.
type Typecaster struct {
	mode CastMode
}

// NewTypecaster returns a Typecaster for the given mode. An empty mode
// defaults to CastNone (passthrough).
func NewTypecaster(mode CastMode) *Typecaster {
	if mode == "" {
		mode = CastNone
	}
	return &Typecaster{mode: mode}
}

// Apply iterates over env pairs and normalises each value according to the
// configured mode. Malformed pairs are passed through unchanged.
func (t *Typecaster) Apply(pairs []string) ([]string, error) {
	if t.mode == CastNone {
		return pairs, nil
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
		norm, err := t.cast(val)
		if err != nil {
			return nil, fmt.Errorf("typecast: key %q: %w", key, err)
		}
		out = append(out, key+"="+norm)
	}
	return out, nil
}

func (t *Typecaster) cast(val string) (string, error) {
	switch t.mode {
	case CastBool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return val, nil // leave non-bool values unchanged
		}
		return strconv.FormatBool(b), nil
	case CastNumeric:
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return strconv.FormatInt(i, 10), nil
		}
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return strconv.FormatFloat(f, 'f', -1, 64), nil
		}
		return val, nil
	case CastAuto:
		if b, err := strconv.ParseBool(val); err == nil {
			return strconv.FormatBool(b), nil
		}
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return strconv.FormatInt(i, 10), nil
		}
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return strconv.FormatFloat(f, 'f', -1, 64), nil
		}
		return val, nil
	default:
		return val, nil
	}
}
