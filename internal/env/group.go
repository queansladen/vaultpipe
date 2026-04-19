package env

import "strings"

// Grouper partitions a slice of KEY=VALUE pairs into named buckets
// based on a prefix map. Keys that match no prefix land in the
// "" (default) bucket.
type Grouper struct {
	prefixes map[string]string // prefix -> bucket name
}

// NewGrouper creates a Grouper from a map of prefix->bucket entries.
func NewGrouper(prefixes map[string]string) *Grouper {
	return &Grouper{prefixes: prefixes}
}

// Group partitions pairs into buckets. Each pair is placed in the
// first matching bucket (iteration order is non-deterministic for
// maps, so callers should use unambiguous prefixes). Pairs that
// match no prefix are placed in the default bucket ("").
func (g *Grouper) Group(pairs []string) map[string][]string {
	buckets := make(map[string][]string)

loop:
	for _, pair := range pairs {
		key := pair
		if idx := strings.IndexByte(pair, '='); idx >= 0 {
			key = pair[:idx]
		}
		for prefix, bucket := range g.prefixes {
			if strings.HasPrefix(key, prefix) {
				buckets[bucket] = append(buckets[bucket], pair)
				continue loop
			}
		}
		// default bucket
		buckets[""] = append(buckets[""], pair)
	}
	return buckets
}

// Flatten merges buckets back into a single slice in bucket-name
// alphabetical order, with the default bucket last.
func Flatten(buckets map[string][]string) []string {
	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		if k != "" {
			keys = append(keys, k)
		}
	}
	// simple insertion sort – bucket count is expected to be small
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	keys = append(keys, "") // default bucket last

	var out []string
	for _, k := range keys {
		out = append(out, buckets[k]...)
	}
	return out
}
