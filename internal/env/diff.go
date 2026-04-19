package env

// ChangeType describes how an environment variable changed.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
	ChangeUpdated ChangeType = "updated"
)

// Change represents a single environment variable change.
type Change struct {
	Key    string
	Old    string
	New    string
	Change ChangeType
}

// Diff compares two slices of KEY=VALUE pairs and returns the changes.
// before and after are both in KEY=VALUE format.
func Diff(before, after []string) []Change {
	prev := toMap(before)
	next := toMap(after)

	var changes []Change

	for k, v := range next {
		if old, ok := prev[k]; !ok {
			changes = append(changes, Change{Key: k, Old: "", New: v, Change: ChangeAdded})
		} else if old != v {
			changes = append(changes, Change{Key: k, Old: old, New: v, Change: ChangeUpdated})
		}
	}

	for k, v := range prev {
		if _, ok := next[k]; !ok {
			changes = append(changes, Change{Key: k, Old: v, New: "", Change: ChangeRemoved})
		}
	}

	return changes
}

func toMap(pairs []string) map[string]string {
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		for i := 0; i < len(p); i++ {
			if p[i] == '=' {
				m[p[:i]] = p[i+1:]
				break
			}
		}
	}
	return m
}
