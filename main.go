package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello Audigent")

	ttlCache := NewCache(2 * time.Second)

	ttlCache.Set([]byte("entry-1"), []byte("value #1"), 1*time.Second)
	ttlCache.Set([]byte("entry-2"), []byte("value #2"), 3*time.Second)
	ttlCache.Set([]byte("entry-3"), []byte("value #3"), 5*time.Second)
	ttlCache.Set([]byte("entry-4"), []byte("value #4"), 10*time.Minute)

	val, ttl := ttlCache.Get([]byte("entry-1"))
	fmt.Printf("Value %+v  - TTL: %+v\n\n", val, ttl)
	time.Sleep(3 * time.Second)
	val, ttl = ttlCache.Get([]byte("entry-1"))
	fmt.Printf("Value %+v  - TTL: %+v\n\n", val, ttl)
	val, ttl = ttlCache.Get([]byte("entry-3"))
	fmt.Printf("Value %+v  - TTL: %+v\n\n", val, ttl)
	time.Sleep(3 * time.Second)
	val, ttl = ttlCache.Get([]byte("entry-3"))
	fmt.Printf("Value %+v  - TTL: %+v\n\n", val, ttl)
	val, ttl = ttlCache.Get([]byte("entry-4"))
	fmt.Printf("Value %+v  - TTL: %+v\n\n", val, ttl)
}
