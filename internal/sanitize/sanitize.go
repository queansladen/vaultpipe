package sanitize

import (
	"errors"
	"strings"
	"unicode"
)

// ErrInvalidKey is returned when an env var key fails validation.
var ErrInvalidKey = errors.New("invalid environment variable key")

// Key validates and normalises an environment variable key.
// It uppercases the key and replaces hyphens and spaces with underscores.
// Returns ErrInvalidKey if the key is empty, starts with a digit, or contains
// characters outside [A-Z0-9_].
func Key(raw string) (string, error) {
	if raw == "" {
		return "", ErrInvalidKey
	}

	upper := strings.ToUpper(strings.NewReplacer("-", "_", " ", "_").Replace(raw))

	for i, r := range upper {
		if i == 0 && unicode.IsDigit(r) {
			return "", ErrInvalidKey
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return "", ErrInvalidKey
		}
	}

	return upper, nil
}

// Value trims leading and trailing whitespace from a secret value.
// It does not modify the interior of the value to avoid corrupting secrets
// such as PEM blocks or JSON payloads.
func Value(v string) string {
	return strings.TrimSpace(v)
}
