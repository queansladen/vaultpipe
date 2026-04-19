package config

// ExpandSecretRefs walks the secret mappings and expands any
// environment variable references found in Path or Key fields
// using the provided lookup function. This allows config values
// like path: "secret/${ENV}/db" to be resolved at runtime.
func ExpandSecretRefs(mappings []SecretMapping, lookup func(string) string) []SecretMapping {
	expanded := make([]SecretMapping, len(mappings))
	for i, m := range mappings {
		expanded[i] = SecretMapping{
			Env:  m.Env,
			Path: expandString(m.Path, lookup),
			Key:  expandString(m.Key, lookup),
		}
	}
	return expanded
}

func expandString(s string, lookup func(string) string) string {
	if lookup == nil {
		return s
	}
	out := make([]byte, 0, len(s))
	i := 0
	for i < len(s) {
		if s[i] == '$' && i+1 < len(s) {
			if s[i+1] == '{' {
				end := i + 2
				for end < len(s) && s[end] != '}' {
					end++
				}
				key := s[i+2 : end]
				out = append(out, lookup(key)...)
				i = end + 1
				continue
			}
		}
		out = append(out, s[i])
		i++
	}
	return string(out)
}
