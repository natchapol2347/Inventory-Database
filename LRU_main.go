package main

import (

	// "errors"
	"fmt"
	// "log"
	// "mime/multipart"
	"sync"
	// "time"
)

type Cache struct {
	capacity uint16
	items    map[string]*cacheItem
	mu       sync.Mutex
	head     *cacheItem
	tail     *cacheItem
}

type cacheItem struct {
	quantity uint32
	next     *cacheItem
	prev	 *cacheItem
}

func newCache(c uint16) *Cache {
	return &Cache{
		capacity: c,
		items:    make(map[string]*cacheItem),
		mu:       sync.Mutex{},
		head:     nil,
		tail: 	  nil,
	}
}

func newItem(key uint32, value uint32) *cacheItem{
	return &cacheItem{
		quantity: value,
		next: nil,
		prev: nil,
	}
}
func (c *Cache) insert_tail(newItem *cacheItem){
	if(c.tail == nil){
		c.tail = newItem
		c.head = newItem

	}else{
		current := c.head
		for(current.prev.prev != nil){
			temp := current.prev
			current.prev = current.prev.prev
			current = temp	 
		}

		c.tail = newItem
	}

}
func (c *Cache) printCache(){
	for key, item := range c.items{
		fmt.Printf("%d : %d", key, item.quantity)
		fmt.Println("hey boyt")
	}
}
func main(){
	x := newCache(3)
	y := newItem(0123,5000)
	x.insert_tail(y)
	x.printCache()
}