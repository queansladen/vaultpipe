package env

import "fmt"

// Chunker splits a slice of KEY=VALUE pairs into fixed-size batches.
// This is useful when spawning processes that have a hard limit on the
// number of environment variables they can receive at once.
type Chunker struct {
	size int
}

// NewChunker returns a Chunker that splits slices into batches of size n.
// If n is <= 0 the entire slice is returned as a single chunk.
func NewChunker(size int) (*Chunker, error) {
	if size < 0 {
		return nil, fmt.Errorf("env/chunk: size must be >= 0, got %d", size)
	}
	return &Chunker{size: size}, nil
}

// Chunk partitions pairs into consecutive slices of at most c.size entries.
// When size is 0 the input is returned as a single chunk without copying.
func (c *Chunker) Chunk(pairs []string) [][]string {
	if len(pairs) == 0 {
		return nil
	}
	if c.size == 0 {
		return [][]string{pairs}
	}

	var chunks [][]string
	for start := 0; start < len(pairs); start += c.size {
		end := start + c.size
		if end > len(pairs) {
			end = len(pairs)
		}
		chunk := make([]string, end-start)
		copy(chunk, pairs[start:end])
		chunks = append(chunks, chunk)
	}
	return chunks
}

// Merge is the inverse of Chunk: it concatenates all chunks into a flat slice.
func Merge(chunks [][]string) []string {
	var total int
	for _, c := range chunks {
		total += len(c)
	}
	out := make([]string, 0, total)
	for _, c := range chunks {
		out = append(out, c...)
	}
	return out
}
