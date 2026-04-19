package env

// Dedupe removes duplicate environment variable entries, keeping the last
// occurrence of each key. This mirrors the behaviour of most Unix shells
// when duplicate keys are present in an environment slice.
func Dedupe(env []string) []string {
	seen := make(map[string]int, len(env))
	order := make([]string, 0, len(env))

	for _, pair := range env {
		key, _ := splitPair(pair)
		if idx, exists := seen[key]; exists {
			// Overwrite the previous entry in-place via order slice.
			order[idx] = pair
		} else {
			seen[key] = len(order)
			order = append(order, pair)
		}
	}

	return order
}

// DedupeMap converts a deduplicated env slice into a map for O(1) lookups.
func DedupeMap(env []string) map[string]string {
	result := make(map[string]string, len(env))
	for _, pair := range Dedupe(env) {
		k, v := splitPair(pair)
		result[k] = v
	}
	return result
}
