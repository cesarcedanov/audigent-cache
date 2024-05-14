package main

import "time"

type CacheEntry struct {
	Value          []byte
	ExpirationTime time.Time
}

// IsExpired will compare the ExpirationTime with the current time.
func (ce CacheEntry) IsExpired() bool {
	// return true if the expiration time already passed
	return ce.ExpirationTime.Before(time.Now())
}
