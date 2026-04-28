package env

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// HashMode controls which algorithm is used to hash environment values.
type HashMode string

const (
	HashModeNone   HashMode = "none"
	HashModeMD5    HashMode = "md5"
	HashModeSHA256 HashMode = "sha256"
)

// Hasher replaces environment variable values with their hashed equivalents.
type Hasher struct {
	mode HashMode
}

// NewHasher returns a Hasher configured with the given mode.
// An empty or unrecognised mode defaults to HashModeNone.
func NewHasher(mode HashMode) *Hasher {
	switch mode {
	case HashModeMD5, HashModeSHA256:
		return &Hasher{mode: mode}
	default:
		return &Hasher{mode: HashModeNone}
	}
}

// Apply hashes the value of each KEY=VALUE pair according to the configured mode.
// Malformed entries (no "=" separator) are passed through unchanged.
func (h *Hasher) Apply(pairs []string) ([]string, error) {
	if h.mode == HashModeNone {
		return pairs, nil
	}
	out := make([]string, len(pairs))
	for i, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out[i] = p
			continue
		}
		key := p[:idx]
		val := p[idx+1:]
		out[i] = fmt.Sprintf("%s=%s", key, h.hash(val))
	}
	return out, nil
}

func (h *Hasher) hash(s string) string {
	switch h.mode {
	case HashModeMD5:
		sum := md5.Sum([]byte(s))
		return hex.EncodeToString(sum[:])
	case HashModeSHA256:
		sum := sha256.Sum256([]byte(s))
		return hex.EncodeToString(sum[:])
	default:
		return s
	}
}
