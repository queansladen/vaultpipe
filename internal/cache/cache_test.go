package cache_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/cache"
)

func sampleData() map[string]string {
	return map[string]string{"username": "admin", "password": "s3cr3t"}
}

func TestGet_Miss(t *testing.T) {
	c := cache.New(5 * time.Minute)
	_, ok := c.Get("secret/missing")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestSet_ThenGet(t *testing.T) {
	c := cache.New(5 * time.Minute)
	c.Set("secret/app", sampleData())
	data, ok := c.Get("secret/app")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if data["username"] != "admin" {
		t.Errorf("unexpected value: %s", data["username"])
	}
}

func TestGet_ExpiredEntry(t *testing.T) {
	c := cache.New(1 * time.Millisecond)
	c.Set("secret/app", sampleData())
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestGet_ZeroTTL_NeverExpires(t *testing.T) {
	c := cache.New(0)
	c.Set("secret/app", sampleData())
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("secret/app")
	if !ok {
		t.Fatal("expected zero-TTL entry to never expire")
	}
}

func TestInvalidate(t *testing.T) {
	c := cache.New(5 * time.Minute)
	c.Set("secret/app", sampleData())
	c.Invalidate("secret/app")
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected entry to be invalidated")
	}
}

func TestLen(t *testing.T) {
	c := cache.New(5 * time.Minute)
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
	c.Set("secret/a", sampleData())
	c.Set("secret/b", sampleData())
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}
