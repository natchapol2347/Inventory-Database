package main

import (

	// "errors"
	"fmt"
	// "math/big"
	// "mime/multipart"
	
	"sync"
	"time"

	
)

//Cache codes
type Cache struct {
	capacity int
	size int
	items    sync.Map
	mu       sync.Mutex
	lookUp   sync.RWMutex
	head     *cacheItem
	tail     *cacheItem
}

type cacheItem struct {
	name     string
	serial   int
	quantity int
	next     *cacheItem
	prev	 *cacheItem
	last_promoted time.Time
}

func newCache(c int) *Cache {
	return &Cache{
		capacity: c,
		size: 0,
		items:    sync.Map{},
		mu:       sync.Mutex{},
		lookUp:	  sync.RWMutex{},
		head:     nil,
		tail: 	  nil,
	}
}

func newItemNode(in_name string, key int, value int) *cacheItem{
	return &cacheItem{
		name: in_name,
		serial: key,
		quantity: value,
		next: nil,
		prev: nil,
		last_promoted: time.Time{},
	}
}

func (c *cacheItem) addi(x int){
	c.quantity += x
}

func main(){
	cache := newCache(6)
	
	cache.items.Store(3, newItemNode("babe",3, 60))
	fmt.Println(cache.items.Load(3))
	var x *cacheItem
	x, _ = (cache.items.Load(3)).cacheItem

	// res, _ := cache.items.Load(3).quantity 
	// res.addi(30)
	// fmt.Println(cache.items.Load(3))

}