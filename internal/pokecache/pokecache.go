package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time //represents when entry was created
	val []byte //data we're caching
}

type Cache struct {
	cache map[string]cacheEntry
	mu *sync.Mutex //for mutex?
}

func NewCache(interval time.Duration) {
	// creates a new cache with conviguratble "interval"



}

func (c *Cache) Add(key string, val []byte) {
	mu.Lock()
	defer mu.Unlock()


}

func (c *Cache) Get(key string) []byte, bool {
	mu.Lock()
	defer mu.Unlock()

}

func (c *Cache) reapLoop() {
	//do i want to lock this? i think it's a helper app for elsewhere, so it might already be locked to 
	//where the call for this function is coming from
	mu.Lock()
	defer mu.Unlock()
	//time.Ticker for this?
}