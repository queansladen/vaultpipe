package env

import (
	"fmt"
	"os"
	"strings"
)

// Builder constructs an environment slice suitable for exec.Cmd.Env.
type Builder struct {
	base    []string
	overlay map[string]string
}

// NewBuilder creates a Builder seeded with the current process environment.
func NewBuilder() *Builder {
	return &Builder{
		base:    os.Environ(),
		overlay: make(map[string]string),
	}
}

// NewBuilderFromBase creates a Builder seeded with the provided base slice.
func NewBuilderFromBase(base []string) *Builder {
	return &Builder{
		base:    base,
		overlay: make(map[string]string),
	}
}

// Set adds or overrides an environment variable in the overlay.
func (b *Builder) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("env: key must not be empty")
	}
	if strings.ContainsRune(key, '=') {
		return fmt.Errorf("env: key %q must not contain '='" , key)
	}
	b.overlay[key] = value
	return nil
}

// SetAll calls Set for each entry in the provided map.
func (b *Builder) SetAll(secrets map[string]string) error {
	for k, v := range secrets {
		if err := b.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

// Build returns the merged environment slice. Overlay values take precedence
// over base values for duplicate keys.
func (b *Builder) Build() []string {
	result := make(map[string]string, len(b.base)+len(b.overlay))

	for _, entry := range b.base {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}

	for k, v := range b.overlay {
		result[k] = v
	}

	out := make([]string, 0, len(result))
	for k, v := range result {
		out = append(out, k+"="+v)
	}
	return out
}

// Keys returns the list of overlay keys that were set.
func (b *Builder) Keys() []string {
	keys := make([]string, 0, len(b.overlay))
	for k := range b.overlay {
		keys = append(keys, k)
	}
	return keys
}
