package cache

import (
	"sync"
	"time"
)

// Entry holds a cached secret value with an expiry.
type Entry struct {
	Data      map[string]string
	FetchedAt time.Time
	TTL       time.Duration
}

// Expired returns true if the entry is past its TTL.
func (e Entry) Expired() bool {
	if e.TTL == 0 {
		return false
	}
	return time.Since(e.FetchedAt) > e.TTL
}

// Cache is a simple in-memory store for secret data keyed by Vault path.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
	ttl     time.Duration
}

// New creates a Cache with the given default TTL. A TTL of 0 disables expiry.
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

// Get retrieves secret data for path. Returns (data, true) on a valid hit.
func (c *Cache) Get(path string) (map[string]string, bool) {
	c.mu.RLock()
	entry, ok := c.entries[path]
	c.mu.RUnlock()
	if !ok || entry.Expired() {
		return nil, false
	}
	return entry.Data, true
}

// Set stores secret data for path using the cache's default TTL.
func (c *Cache) Set(path string, data map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[path] = Entry{
		Data:      data,
		FetchedAt: time.Now(),
		TTL:       c.ttl,
	}
}

// Invalidate removes the entry for path.
func (c *Cache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, path)
}

// Len returns the number of entries currently in the cache.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
