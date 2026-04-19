package env

// PassthroughMode controls which environment variables from the host
// are forwarded to the child process.
type PassthroughMode int

const (
	// PassthroughAll forwards every variable from the host environment.
	PassthroughAll PassthroughMode = iota
	// PassthroughNone starts the child with an empty environment.
	PassthroughNone
	// PassthroughFiltered forwards only variables matching the supplied filter.
	PassthroughFiltered
)

// Passthrough builds an env slice for the child process by applying mode
// logic on top of a base snapshot and an overlay of resolved secrets.
//
// overlay values always win over inherited host variables.
func Passthrough(mode PassthroughMode, host *Snapshot, filter *Filter, overlay map[string]string) []string {
	var base map[string]string

	switch mode {
	case PassthroughAll:
		base = host.All()
	case PassthroughNone:
		base = make(map[string]string)
	case PassthroughFiltered:
		if filter == nil {
			base = host.All()
		} else {
			base = filteredMap(host, filter)
		}
	}

	return MergeEnv(mapToSlice(base), mapToSlice(overlay))
}

func filteredMap(s *Snapshot, f *Filter) map[string]string {
	out := make(map[string]string)
	for _, pair := range f.Apply(s.Pairs()) {
		k, v, ok := splitPair(pair)
		if ok {
			out[k] = v
		}
	}
	return out
}

func mapToSlice(m map[string]string) []string {
	slice := make([]string, 0, len(m))
	for k, v := range m {
		slice = append(slice, k+"="+v)
	}
	return slice
}

func splitPair(pair string) (string, string, bool) {
	for i := 0; i < len(pair); i++ {
		if pair[i] == '=' {
			return pair[:i], pair[i+1:], true
		}
	}
	return "", "", false
}
