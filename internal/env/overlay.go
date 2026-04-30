package env

// Overlay merges two slices of KEY=VALUE pairs, with entries from the overlay
// slice taking precedence over those in the base slice when keys collide.
//
// The result preserves the relative ordering of the base slice, replacing
// in-place any key that appears in the overlay, then appending any overlay
// keys that were not present in the base.
//
// This is distinct from Merge (which accepts a strategy) and Join (which
// concatenates values): Overlay always lets the overlay win, matching the
// semantics expected when injecting Vault secrets on top of an inherited
// process environment.
package env

import "strings"

// Overlayer combines a base environment with an overlay environment.
type Overlayer struct {
	ignoreCase bool
}

// OverlayOption configures an Overlayer.
type OverlayOption func(*Overlayer)

// WithCaseInsensitiveKeys makes key comparison case-insensitive.
// This is useful on platforms (e.g. Windows) where environment variable names
// are conventionally treated as case-insensitive.
func WithCaseInsensitiveKeys() OverlayOption {
	return func(o *Overlayer) {
		o.ignoreCase = true
	}
}

// NewOverlayer creates an Overlayer with the supplied options.
func NewOverlayer(opts ...OverlayOption) *Overlayer {
	o := &Overlayer{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Apply merges overlay on top of base and returns the resulting slice.
// Entries in overlay whose keys already exist in base replace those entries
// in-place; remaining overlay entries are appended in their original order.
func (o *Overlayer) Apply(base, overlay []string) []string {
	if len(overlay) == 0 {
		out := make([]string, len(base))
		copy(out, base)
		return out
	}

	// Build a lookup from normalised key → overlay value for fast access.
	overlayMap := make(map[string]string, len(overlay))
	for _, pair := range overlay {
		k, v, ok := splitPairOverlay(pair)
		if !ok {
			continue
		}
		overlayMap[o.normalise(k)] = v
	}

	// Walk base, replacing any key found in the overlay.
	result := make([]string, 0, len(base)+len(overlay))
	seen := make(map[string]struct{}, len(base))

	for _, pair := range base {
		k, _, ok := splitPairOverlay(pair)
		if !ok {
			// Preserve malformed entries unchanged.
			result = append(result, pair)
			continue
		}
		norm := o.normalise(k)
		seen[norm] = struct{}{}

		if v, exists := overlayMap[norm]; exists {
			result = append(result, k+"="+v)
		} else {
			result = append(result, pair)
		}
	}

	// Append overlay entries whose keys were not present in base.
	for _, pair := range overlay {
		k, _, ok := splitPairOverlay(pair)
		if !ok {
			continue
		}
		if _, exists := seen[o.normalise(k)]; !exists {
			result = append(result, pair)
		}
	}

	return result
}

// normalise returns the key in the form used for map lookups.
func (o *Overlayer) normalise(key string) string {
	if o.ignoreCase {
		return strings.ToUpper(key)
	}
	return key
}

// splitPairOverlay splits a KEY=VALUE string into its components.
// It returns (key, value, true) on success or ("", "", false) for malformed input.
func splitPairOverlay(pair string) (string, string, bool) {
	idx := strings.IndexByte(pair, '=')
	if idx < 0 {
		return "", "", false
	}
	return pair[:idx], pair[idx+1:], true
}
