package main

import (
	"reflect"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	defaultTTL := 5 * time.Second
	testCache := NewCache(defaultTTL)

	t.Run("ok", func(t *testing.T) {
		defer testCache.EmptyEntries()
		testCache.Set([]byte("foo"), []byte("bar"), 2*time.Second)
		value, ttl := testCache.Get([]byte("foo"))
		if !reflect.DeepEqual(value, []byte("bar")) {
			t.Fatalf("expected: %v, got %v", []byte("bar"), value)
		}
		if ttl < 0*time.Second {
			t.Fatalf("%v should be higher than 0", ttl)
		}
	})
	t.Run("empty storage", func(t *testing.T) {
		// Clean up cache
		testCache.EmptyEntries()
		if testCache.TotalEntries() != 0 {
			t.Fatalf("Cache should be Empty, but got %v entries on it", len(testCache.Entries))
		}
	})
	t.Run("set value with higher ttl than routine", func(t *testing.T) {
		defer testCache.EmptyEntries()
		testCache.Set([]byte("foo"), []byte("bar"), defaultTTL*2)
		time.Sleep(defaultTTL)
		value, _ := testCache.Get([]byte("foo"))
		if !reflect.DeepEqual(value, []byte("bar")) {
			t.Fatalf("expected: %v, got %v", []byte("bar"), value)
		}
		if testCache.TotalEntries() != 1 {
			t.Fatalf("Cache should contains just 1 entry but got %v entries on it", testCache.TotalEntries())
		}
	})

	t.Run("set value with lower ttl than routine", func(t *testing.T) {
		defer testCache.EmptyEntries()
		testCache.Set([]byte("foo"), []byte("bar"), 1*time.Second)

		time.Sleep(defaultTTL)
		value, ttl := testCache.Get([]byte("foo"))
		if !reflect.DeepEqual(value, []byte{}) {
			t.Fatalf("expected: %v, got %v", []byte{}, value)
		}

		if testCache.TotalEntries() != 0 {
			t.Fatalf("Cache should be Empty but got %v entries on it. \n\n %+v \n\n\n %+v", testCache.TotalEntries(), ttl, testCache.Entries)
		}
	})

	t.Run("set multiple values with higher/lower ttl than routine", func(t *testing.T) {
		defer testCache.EmptyEntries()
		testCases := []struct {
			key                            []byte
			value                          []byte
			expiration                     time.Duration
			expectedValue1, expectedValue2 []byte
		}{
			{
				key:            []byte("foo"),
				value:          []byte("bar"),
				expiration:     2 * time.Second,
				expectedValue1: []byte("bar"),
				expectedValue2: []byte{},
			},
			{
				key:            []byte("hello"),
				value:          []byte("audigent"),
				expiration:     10 * time.Second,
				expectedValue1: []byte("audigent"),
				expectedValue2: []byte("audigent"),
			},
		}
		// Set entries in Cache
		for _, tc := range testCases {
			testCache.Set(tc.key, tc.value, tc.expiration)
			value, _ := testCache.Get(tc.key)
			if !reflect.DeepEqual(value, tc.expectedValue1) {
				t.Fatalf("expected: %v, got %v", tc.expectedValue1, value)
			}
		}

		if testCache.TotalEntries() != 2 {
			t.Fatalf("Cache should contains just 2 entry but got %v entries on it", testCache.TotalEntries())
		}
		// Trigger Routine
		time.Sleep(defaultTTL)

		// Get entries in Cache after Active strategy got trigger
		for _, tc := range testCases {
			value, _ := testCache.Get(tc.key)
			if !reflect.DeepEqual(value, tc.expectedValue2) {
				t.Fatalf("expected: %v, got %v", tc.expectedValue2, value)
			}
		}

		if testCache.TotalEntries() != 1 {
			t.Fatalf("Cache should contains just 1 entry but got %v entries on it", testCache.TotalEntries())
		}
	})

}

var benchmarkCases = []struct {
	name       string
	key, val   []byte
	expiration time.Duration
}{
	{
		name:       "Quick TTL",
		key:        []byte("audigent"),
		val:        []byte("audigent"),
		expiration: 100 * time.Millisecond,
	}, {
		name:       "Long TTL",
		key:        []byte("foo"),
		val:        []byte("bar"),
		expiration: 2000 * time.Millisecond,
	},
}

func BenchmarkSet(b *testing.B) {
	bCache := NewCache(1000 * time.Millisecond)
	for _, bc := range benchmarkCases {
		b.Run(bc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				bCache.Set(bc.key, bc.val, bc.expiration)
			}
		})
	}

}

func BenchmarkGet(b *testing.B) {
	bCache := NewCache(1000 * time.Millisecond)
	for _, bc := range benchmarkCases {
		b.Run(bc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				bCache.Set(bc.key, bc.val, bc.expiration)
				bCache.Get(bc.key)
			}
		})
	}

}
