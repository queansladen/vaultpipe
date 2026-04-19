package env

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError holds all validation failures for an environment slice.
type ValidationError struct {
	Errs []string
}

func (e *ValidationError) Error() string {
	return "env validation failed:\n  " + strings.Join(e.Errs, "\n  ")
}

// Validator checks environment variable slices for common problems.
type Validator struct {
	allowEmpty bool
}

// NewValidator returns a Validator. When allowEmpty is false, entries with
// empty values are reported as warnings but not errors.
func NewValidator(allowEmpty bool) *Validator {
	return &Validator{allowEmpty: allowEmpty}
}

// Validate checks each entry in pairs (KEY=VALUE format) and returns a
// *ValidationError if any hard violations are found.
func (v *Validator) Validate(pairs []string) error {
	var errs []string
	seen := make(map[string]int)

	for i, pair := range pairs {
		idx := strings.IndexByte(pair, '=')
		if idx < 0 {
			errs = append(errs, fmt.Sprintf("entry %d: missing '=' separator: %q", i, pair))
			continue
		}
		key := pair[:idx]
		value := pair[idx+1:]

		if key == "" {
			errs = append(errs, fmt.Sprintf("entry %d: empty key", i))
			continue
		}

		if !v.allowEmpty && value == "" {
			errs = append(errs, fmt.Sprintf("key %q: empty value", key))
		}

		if prev, ok := seen[key]; ok {
			errs = append(errs, fmt.Sprintf("key %q: duplicate (first at entry %d, repeated at entry %d)", key, prev, i))
		} else {
			seen[key] = i
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return &ValidationError{Errs: errs}
}

// IsValidationError reports whether err is a *ValidationError.
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}
