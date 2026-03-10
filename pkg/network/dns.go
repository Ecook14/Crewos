package network

import (
	"sync"
	"time"
)

type dnsEntry struct {
	IP        string
	ExpiresAt time.Time
}

// DNSCache provides a lightweight in-memory cache for LLM API lookups.
type DNSCache struct {
	entries map[string]dnsEntry
	mu      sync.RWMutex
	ttl     time.Duration
}

// NewDNSCache initializes the caching layer.
func NewDNSCache(ttl time.Duration) *DNSCache {
	return &DNSCache{
		entries: make(map[string]dnsEntry),
		ttl:     ttl,
	}
}

// Get retrieves a cached IP if it exists and hasn't expired.
func (c *DNSCache) Get(host string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[host]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return "", false
	}
	return entry.IP, true
}

// Set stores a resolution in the cache.
func (c *DNSCache) Set(host, ip string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[host] = dnsEntry{
		IP:        ip,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}
