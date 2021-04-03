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
	capacity int
	size int
	items    map[int]*cacheItem
	mu       sync.Mutex
	head     *cacheItem
	tail     *cacheItem
}

type cacheItem struct {
	serial   int
	quantity int
	next     *cacheItem
	prev	 *cacheItem
}

func newCache(c int) *Cache {
	return &Cache{
		capacity: c,
		size: 0,
		items:    make(map[int]*cacheItem),
		mu:       sync.Mutex{},
		head:     nil,
		tail: 	  nil,
	}
}

func newItemNode(key int, value int) *cacheItem{
	return &cacheItem{
		serial: key,
		quantity: value,
		next: nil,
		prev: nil,
	}
}
func (c *Cache) insert_tail(key int, value int) *cacheItem{
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
		c.head = c.head.prev
		c.head.next = nil
	}
}
// func (c *Cache) printCache(){
// 	for key, item := range c.items{
// 		fmt.Printf("%d : %d", key, item.quantity)
// 		fmt.Println("hey boyt")
// 	}
// }

func (c *Cache) get(key int, load int) int{
	if _, ok := c.items[key]; ok{
		value := c.items[key].quantity
		if(value - load < 0){
			fmt.Println("not enough in stock")
			return 0 //error handling might change in the future
			
		}
		c.moveToFront(c.items[key])
		c.items[key].quantity -= load 
		fmt.Println("The quantity remaining is", value - load )
		return value
	}else{
		//if there's no key
		fmt.Println("there's no key", key, "yet")
		return 0 //error handiling might change in the future
	}
}

func (c *Cache) put(key int, load int) {
	if _, ok := c.items[key]; ok {
		c.items[key].quantity += load
		c.moveToFront(c.items[key])
		return
	}

	if c.size == c.capacity {
		delKey := c.head.serial
		c.pop()
		c.size--
		delete(c.items, delKey)
	}
	page := c.insert_tail(key, load)
	c.size++
	c.items[key] = page
}

func main() {
	cache := newCache(2)
	cache.put(2, 2)
	cache.get(2,100)
	cache.get(1,100)
	cache.put(1, 1)
	cache.put(1, 5)
	cache.get(1, 5)
	cache.get(2,25)
	cache.put(8, 8)
	cache.get(1, 3)
	cache.get(8, 5)
}