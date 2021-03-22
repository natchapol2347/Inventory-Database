// Package gCache implements an LRU cache in golang.
//   Set: Add the item both in queue and HashMap. If they capacity is full,
//        it removes the least recently used element.
//
//   Get: Returns the item requested via Key. On querying the item it comes
//        to forward of the queue
package main

import (
	"errors"
	"sync"
	"time"
	"container/list"
	"fmt"

)

// Cache is an object which will hold items, it is the cache of these items.
type Cache struct {
	capacity     int
	items        map[string]*cacheItem
	mu           sync.Mutex
	timesEvicted int
	order        *list.List
}

type cacheItem struct {
	value        string
	lastTimeUsed int64
}

// Create a new cache object.
func New(c int) *Cache {
	return &Cache{
		capacity: c,
		items:    make(map[string]*cacheItem),
		mu:       sync.Mutex{},
		order:    list.New(),
	}
}

// Set a key into the cache, remove the last used key if capacity has been met.
func (c *Cache) Set(key string, val string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Search for the key in map, if the key isn't there
	// add it, no action if the key already exists.
	if _, ok := c.items[key]; !ok {
		//record time accessed
		now := time.Now().UnixNano()
		n := c.capacity
		if len(c.items) == n { // If more than cache capacity- evict
			
			delete(c.items, del)
			
		}

		// Add the new element to the cache.
		c.items[key] = &cacheItem{
			value:        val,
			lastTimeUsed: now,
		}
	}
}

// Get a key from the cache and update that key's lastTimeUsed.
func (c *Cache) Get(k string) (string, error) {
	//Search the key in map
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.items[k]; ok {
		
		return v.value, nil
	}
	else{
		evict()
	}
}

func main(){
	//user adds 5 items to cache
	itemCache := New(5)
	itemCache.Set("banana","234")
	itemCache.Set("doll","235")
	itemCache.Set("steak","236")
	itemCache.Set("ball","237")
	itemCache.Set("tv","239")
	itemCache.Set("banana","234")
	//call Get() everytime user imports or exports
	fmt.Println(itemCache.Get("banana"))


}