package env

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewChunker_NegativeSize_ReturnsError(t *testing.T) {
	_, err := NewChunker(-1)
	if err == nil {
		t.Fatal("expected error for negative size")
	}
}

func TestNewChunker_ZeroSize_OK(t *testing.T) {
	_, err := NewChunker(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestChunk_EmptyInput_ReturnsNil(t *testing.T) {
	c, _ := NewChunker(3)
	got := c.Chunk(nil)
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestChunk_ZeroSize_SingleChunk(t *testing.T) {
	c, _ := NewChunker(0)
	pairs := []string{"A=1", "B=2", "C=3"}
	got := c.Chunk(pairs)
	if len(got) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(got))
	}
	if diff := cmp.Diff(pairs, got[0]); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestChunk_EvenSplit(t *testing.T) {
	c, _ := NewChunker(2)
	pairs := []string{"A=1", "B=2", "C=3", "D=4"}
	got := c.Chunk(pairs)
	if len(got) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(got))
	}
	want := [][]string{{"A=1", "B=2"}, {"C=3", "D=4"}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestChunk_UnevenSplit_LastChunkSmaller(t *testing.T) {
	c, _ := NewChunker(2)
	pairs := []string{"A=1", "B=2", "C=3"}
	got := c.Chunk(pairs)
	if len(got) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(got))
	}
	if len(got[1]) != 1 {
		t.Errorf("expected last chunk len 1, got %d", len(got[1]))
	}
}

func TestChunk_DoesNotMutateInput(t *testing.T) {
	c, _ := NewChunker(2)
	pairs := []string{"A=1", "B=2", "C=3"}
	orig := make([]string, len(pairs))
	copy(orig, pairs)
	c.Chunk(pairs)
	if diff := cmp.Diff(orig, pairs); diff != "" {
		t.Errorf("input mutated (-want +got):\n%s", diff)
	}
}

func TestMerge_ReconstitutesOriginal(t *testing.T) {
	pairs := []string{"A=1", "B=2", "C=3", "D=4", "E=5"}
	c, _ := NewChunker(2)
	chunks := c.Chunk(pairs)
	got := Merge(chunks)
	if diff := cmp.Diff(pairs, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestMerge_NilChunks_ReturnsEmpty(t *testing.T) {
	got := Merge(nil)
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}
