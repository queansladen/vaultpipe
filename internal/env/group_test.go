package env

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGroup_EmptyPrefixes_AllDefault(t *testing.T) {
	g := NewGrouper(nil)
	pairs := []string{"FOO=1", "BAR=2"}
	buckets := g.Group(pairs)
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if diff := cmp.Diff(pairs, buckets[""]); diff != "" {
		t.Errorf("default bucket mismatch (-want +got):\n%s", diff)
	}
}

func TestGroup_MatchingPrefix_PlacedInBucket(t *testing.T) {
	g := NewGrouper(map[string]string{"AWS_": "aws", "DB_": "db"})
	pairs := []string{"AWS_REGION=us-east-1", "DB_HOST=localhost", "APP_NAME=vaultpipe"}
	buckets := g.Group(pairs)

	if got := buckets["aws"]; len(got) != 1 || got[0] != "AWS_REGION=us-east-1" {
		t.Errorf("aws bucket: %v", got)
	}
	if got := buckets["db"]; len(got) != 1 || got[0] != "DB_HOST=localhost" {
		t.Errorf("db bucket: %v", got)
	}
	if got := buckets[""]; len(got) != 1 || got[0] != "APP_NAME=vaultpipe" {
		t.Errorf("default bucket: %v", got)
	}
}

func TestGroup_MalformedEntry_UsesFullStringAsKey(t *testing.T) {
	g := NewGrouper(map[string]string{"FOO": "foo"})
	pairs := []string{"FOOBAR"} // no '='
	buckets := g.Group(pairs)
	if got := buckets["foo"]; len(got) != 1 {
		t.Errorf("expected malformed entry in foo bucket, got %v", buckets)
	}
}

func TestFlatten_OrderedBucketsThenDefault(t *testing.T) {
	buckets := map[string][]string{
		"z": {"Z_KEY=1"},
		"a": {"A_KEY=2"},
		"":  {"PLAIN=3"},
	}
	got := Flatten(buckets)
	want := []string{"A_KEY=2", "Z_KEY=1", "PLAIN=3"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Flatten mismatch (-want +got):\n%s", diff)
	}
}

func TestFlatten_EmptyBuckets_ReturnsNil(t *testing.T) {
	got := Flatten(map[string][]string{})
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}
