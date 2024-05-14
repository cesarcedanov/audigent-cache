package main

import (
	"testing"
	"time"
)

func TestCacheEntry(t *testing.T) {

	testCases := []struct {
		name                         string
		ttl                          time.Duration
		expectedBool1, expectedBool2 bool
	}{
		{
			name:          "Long TTL",
			ttl:           10 * time.Second,
			expectedBool1: false,
			expectedBool2: false,
		},
		{
			name:          "Quick TTL",
			ttl:           1 * time.Second,
			expectedBool1: false,
			expectedBool2: true,
		},
	}
	for _, tc := range testCases {
		newCacheEntry := &CacheEntry{
			Value:          []byte(tc.name),
			ExpirationTime: time.Now().Add(tc.ttl),
		}
		if newCacheEntry.IsExpired() != tc.expectedBool1 {
			t.Fatalf("Test '%s' got a unexcepeted TTL on the 1st Check", tc.name)
		}
		time.Sleep(2 * time.Second)

		if newCacheEntry.IsExpired() != tc.expectedBool2 {
			t.Fatalf("Test '%s' got a unexcepeted TTL on the 2nd Check", tc.name)
		}

	}
}
