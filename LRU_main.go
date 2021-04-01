package main

import (

	// "errors"
	// "fmt"
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
	}

}

func main(){
	x := New(3)
	x.insert_tail = 

}