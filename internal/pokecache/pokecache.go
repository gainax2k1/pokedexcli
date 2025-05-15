package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time //represents when entry was created
	val       []byte    //data we're caching
}

type Cache struct {
	cache map[string]cacheEntry
	mu    *sync.Mutex //for mutex?
}

func NewCache(interval time.Duration) *Cache {
	// creates a new cache with conviguratble "interval"
	newC := &Cache{ //do i need to use a pointer? YES
		cache: make(map[string]cacheEntry), //corrected (remember to MAKE new MAPs!)
		mu:    &sync.Mutex{},               // corrected. remember your pointers!
	}
	go newC.reapLoop(interval) //?

	return newC

}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	/* clunky implementation
	var newCacheEntry cacheEntry
	newCacheEntry.val = val
	newCacheEntry.createdAt = time.Now()
	c.cache[key] = newCacheEntry
	*/

	// using "struct literal" - much cleaner
	c.cache[key] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	/*UNSAFEW WAY TO DO IT:
	if c.cache[key].val == nil {
		return []byte{}, false
	}
	return c.cache[key].val, true
	*/

	entry, exists := c.cache[key] // safe way to do it!
	if !exists {
		return []byte{}, false
	}
	return entry.val, true

}

func (c *Cache) reapLoop(interval time.Duration) {

	reapTicker := time.NewTicker(interval) // first, make our ticker with the duration
	for {
		// wait for next tick
		<-reapTicker.C

		// NOW lock the mutex when we're actually accessing the map
		c.mu.Lock()
		now := time.Now() //catch the time we locked at for comparing

		// check all the times on entrys here:
		/*
			- Get the current time
			- Compare each entry's creation time with the current time
			- Remove entries that are older than interval
		*/
		// For each entry in the cache
		for key, entry := range c.cache {
			// If the entry is older than the interval
			if now.Sub(entry.createdAt) > interval {
				// Delete it from the cache
				delete(c.cache, key)
			}
		}

		c.mu.Unlock() // now safe to unlock
	}
}
