package main

import (
	"math/rand"
	"sync"
	"time"
)

// Cache will storage some entries with a specific ttl
type Cache struct {
	// We should use RWMutex instead of just Mutex
	Mx                  sync.RWMutex
	Entries             map[string]*CacheEntry
	ActiveStrategyTimer *time.Timer
}

// NewCache return a new instance with an Active strategy
func NewCache(duration time.Duration) *Cache {
	ttlCache := &Cache{
		ActiveStrategyTimer: time.NewTimer(duration),
		Entries:             make(map[string]*CacheEntry),
	}

	go ttlCache.triggerRoutine()
	return ttlCache
}

// Set will write add/update an entry based on the Key
func (ttlCache *Cache) Set(key string, value []byte, duration time.Duration) {
	// Lock now and Unlock at the end of the scope
	ttlCache.Mx.Lock()
	defer ttlCache.Mx.Unlock()

	expirationTime := time.Now().Add(duration)
	ttlCache.Entries[key] = &CacheEntry{Value: value, ExpirationTime: expirationTime}
}

// Get will return the entry value based on the Key
func (ttlCache *Cache) Get(key string) ([]byte, time.Duration) {
	// Lock now and Unlock at the end of the scope
	ttlCache.Mx.RLock()
	defer ttlCache.Mx.RUnlock()

	entry, exists := ttlCache.Entries[key]
	if !exists || entry.IsExpired() {
		return nil, 0
	}
	// return the value and the remaining time in the cache
	return entry.Value, entry.ExpirationTime.Sub(time.Now())
}

// RemoveExpiredEntry removed all the expired entries
func (ttlCache *Cache) RemoveExpiredEntry() {
	// Lock now and Unlock at the end of the scope
	ttlCache.Mx.Lock()
	defer ttlCache.Mx.Unlock()

	for key, entry := range ttlCache.Entries {
		if entry.IsExpired() {
			delete(ttlCache.Entries, key)
		}
	}
}

// triggerRoutine is needed to handle the Active Strategy
// based on (https://www.pankajtanwar.in/blog/how-redis-expires-keys-a-deep-dive-into-how-ttl-works-internally-in-redis)
func (ttlCache *Cache) triggerRoutine() {
	for {
		select {
		case <-ttlCache.ActiveStrategyTimer.C:
			ttlCache.RemoveExpiredEntry()
			ttlCache.ActiveStrategyTimer.Reset(time.Duration(rand.Intn(100)) * time.Millisecond)
		}
	}
}
