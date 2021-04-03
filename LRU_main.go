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

func newItemNode(key uint32, value uint32) *cacheItem{
	return &cacheItem{
		quantity: value,
		next: nil,
		prev: nil,
	}
}
func (c *Cache) insert_tail(key uint32, value uint32) *cacheItem{
	//make new item from argument
	newItem := newItemNode(key, value)
	if(c.tail == nil && c.head == nil){
		c.tail = newItem
		c.head = newItem

	}else{
		newItem.next = c.tail
		c.tail.prev = newItem
		c.tail = newItem
	}

	return newItem

}

func (c *Cache) moveToFront(node *cacheItem){
	if node == c.tail{
		return 
	}else if node == c.head{
		c.head = c.head.prev
		//last node's next must point to nil
		c.head.next = nil 
	}

	node.next = c.tail
	c.tail.prev = node
	c.tail = node
}

func (c *Cache) pop(){
	if c.head == nil && c.tail == nil{
		return
	}else if c.head == c.tail{
		c.head, c.tail = nil, nil
	}else{
		
	}
}
// func (c *Cache) printCache(){
// 	for key, item := range c.items{
// 		fmt.Printf("%d : %d", key, item.quantity)
// 		fmt.Println("hey boyt")
// 	}
// }

func main(){
	x := newCache(3)
	y := newItem(0123,5000)
	x.insert_tail(y)
	x.printCache()
}