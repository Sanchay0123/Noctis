package routing

import (
	"crypto/rand"
	"sync"
	"time"
)

type PacketCacheKey struct {
	Source   [32]byte
	PacketID [16]byte
}

type PacketCache struct {
	mu       sync.Mutex
	capacity int
	ttl      time.Duration
	entries  map[PacketCacheKey]time.Time
	order    []PacketCacheKey
}

func NewPacketCache(capacity int, ttl time.Duration) *PacketCache {
	return &PacketCache{
		capacity: capacity,
		ttl:      ttl,
		entries:  make(map[PacketCacheKey]time.Time),
		order:    make([]PacketCacheKey, 0, capacity),
	}
}

func (c *PacketCache) IsDuplicateOrAdd(key PacketCacheKey) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	var validOrder []PacketCacheKey
	for _, k := range c.order {
		if ts, ok := c.entries[k]; ok {
			if now.Sub(ts) > c.ttl {
				delete(c.entries, k)
			} else {
				validOrder = append(validOrder, k)
			}
		}
	}
	c.order = validOrder

	if _, exists := c.entries[key]; exists {
		return true
	}

	if len(c.order) >= c.capacity {
		oldest := c.order[0]
		delete(c.entries, oldest)
		c.order = c.order[1:]
	}

	c.entries[key] = now
	c.order = append(c.order, key)
	return false
}

func GeneratePacketID() [16]byte {
	var id [16]byte
	rand.Read(id[:])
	return id
}
