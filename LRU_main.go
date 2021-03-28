package main

import (

	"errors"
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

func New(c uint16) *Cache {
	return &Cache{
		capacity: c,
		items:    make(map[string]*cacheItem),
		mu:       sync.Mutex{},
		head:     nil,
		tail: 	  nil,
	}
}

func (c *Cache) insert_tail(x cacheItem){
	
}
