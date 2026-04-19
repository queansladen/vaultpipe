package env

// Wrapper applies a key and value transformation to each env pair.
// It wraps both the key and the value using the provided functions.
// Malformed entries (no '=' separator) are passed through unchanged.
type Wrapper struct {
	wrapKey   func(string) string
	wrapValue func(string) string
}

// NewWrapper returns a Wrapper that applies wrapKey to keys and wrapValue to values.
// Pass nil for either function to leave that part unchanged.
func NewWrapper(wrapKey, wrapValue func(string) string) *Wrapper {
	if wrapKey == nil {
		wrapKey = func(s string) string { return s }
	}
	if wrapValue == nil {
		wrapValue = func(s string) string { return s }
	}
	return &Wrapper{wrapKey: wrapKey, wrapValue: wrapValue}
}

// Apply transforms each pair in the slice and returns a new slice.
func (w *Wrapper) Apply(pairs []string) []string {
	out := make([]string, 0, len(pairs))
	for _, p := range pairs {
		idx := indexOf(p, '=')
		if idx < 0 {
			out = append(out, p)
			continue
		}
		k := w.wrapKey(p[:idx])
		v := w.wrapValue(p[idx+1:])
		out = append(out, k+"="+v)
	}
	return out
}

// indexOf returns the index of the first occurrence of b in s, or -1.
func indexOf(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}
