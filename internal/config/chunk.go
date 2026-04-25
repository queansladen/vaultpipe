package config

import "github.com/frobware/vaultpipe/internal/env"

// ChunkConfig holds optional chunking parameters from the YAML config.
type ChunkConfig struct {
	// Size is the maximum number of environment variables per chunk.
	// A value of 0 disables chunking (all pairs in a single batch).
	Size int `yaml:"size"`
}

// ResolveChunker constructs an env.Chunker from the supplied config.
// When cfg is nil a zero-size (no-op) chunker is returned so callers
// never need to nil-check the result.
func ResolveChunker(cfg *ChunkConfig) (*env.Chunker, error) {
	size := 0
	if cfg != nil && cfg.Size > 0 {
		size = cfg.Size
	}
	return env.NewChunker(size)
}
