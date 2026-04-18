package redact

import "strings"

// Redactor replaces known secret values with a placeholder in log output.
type Redactor struct {
	secrets []string
	mask    string
}

// New creates a Redactor with the given mask string.
func New(mask string) *Redactor {
	if mask == "" {
		mask = "[REDACTED]"
	}
	return &Redactor{mask: mask}
}

// Add registers one or more secret values to be redacted.
func (r *Redactor) Add(values ...string) {
	for _, v := range values {
		if v != "" {
			r.secrets = append(r.secrets, v)
		}
	}
}

// Redact replaces all registered secret values in s with the mask.
func (r *Redactor) Redact(s string) string {
	for _, secret := range r.secrets {
		s = strings.ReplaceAll(s, secret, r.mask)
	}
	return s
}

// Mask returns the configured mask string.
func (r *Redactor) Mask() string {
	return r.mask
}

// Len returns the number of registered secrets.
func (r *Redactor) Len() int {
	return len(r.secrets)
}
