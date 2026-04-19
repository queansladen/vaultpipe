package env

import "sort"

// Sorter sorts environment variable slices in a stable, deterministic order.
type Sorter struct {
	descending bool
}

// NewSorter returns a Sorter. If descending is true, keys are sorted Z→A.
func NewSorter(descending bool) *Sorter {
	return &Sorter{descending: descending}
}

// Sort returns a new slice with entries sorted by key name.
func (s *Sorter) Sort(pairs []string) []string {
	out := make([]string, len(pairs))
	copy(out, pairs)
	sort.SliceStable(out, func(i, j int) bool {
		ki := keyOf(out[i])
		kj := keyOf(out[j])
		if s.descending {
			return ki > kj
		}
		return ki < kj
	})
	return out
}

// keyOf returns the key portion of a KEY=VALUE pair.
func keyOf(pair string) string {
	for i := 0; i < len(pair); i++ {
		if pair[i] == '=' {
			return pair[:i]
		}
	}
	return pair
}
