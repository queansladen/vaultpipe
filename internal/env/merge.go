package env

// MergeStrategy controls how duplicate keys are resolved when merging
// multiple environment slices together.
type MergeStrategy int

const (
	// StrategyLast keeps the last value seen for a given key.
	StrategyLast MergeStrategy = iota
	// StrategyFirst keeps the first value seen for a given key.
	StrategyFirst
)

// Merger combines multiple KEY=VALUE slices into one, resolving duplicates
// according to the configured strategy.
type Merger struct {
	strategy MergeStrategy
}

// NewMerger returns a Merger using the given strategy.
func NewMerger(strategy MergeStrategy) *Merger {
	return &Merger{strategy: strategy}
}

// Merge combines the provided slices in order. Later slices take precedence
// when StrategyLast is used; earlier slices take precedence with StrategyFirst.
func (m *Merger) Merge(slices ...[]string) []string {
	seen := make(map[string]string)
	order := []string{}

	for _, slice := range slices {
		for _, pair := range slice {
			key, val, ok := splitPair(pair)
			if !ok {
				continue
			}
			if _, exists := seen[key]; !exists {
				order = append(order, key)
			}
			switch m.strategy {
			case StrategyFirst:
				if _, exists := seen[key]; !exists {
					seen[key] = val
				}
			default: // StrategyLast
				seen[key] = val
			}
		}
	}

	out := make([]string, 0, len(order))
	for _, k := range order {
		out = append(out, k+"="+seen[k])
	}
	return out
}
