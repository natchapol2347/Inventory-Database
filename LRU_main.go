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
	head     *cacheItem
	tail     *cacheItem
}

type cacheItem struct {
	name     string
	serial   int
	quantity int
	next     *cacheItem
	prev	 *cacheItem
	mu       sync.Mutex
	

}

func newCache(c int) *Cache {
	return &Cache{
		capacity: c,
		size: 0,
		items:    make(map[int]*cacheItem),
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
		mu :     sync.Mutex{},

	}
}
func (c *Cache) insert_tail(name string, key int, value int) *cacheItem{
	//make new item from argument
	newItem := newItemNode(name,key, value)
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
	}else{
		node.prev.next = node.next
		node.next.prev = node.prev
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
func (c *Cache) printCache(){
	current := c.tail
	var i int;
	for i=0;i<=c.size;i++{
		if(current != nil){		
			fmt.Printf("|name:%s|id:%d|,quantity:%d|size:%d| ->", current.name, current.serial, current.quantity,c.size)
			current = current.next
		}
	}
	fmt.Println("\n")
}


func (c *Cache) get(name string, key int, load int) (int,int){
	if _, ok := c.items[key]; ok{
		value := c.items[key].quantity
		if(value - load < 0){
			// fmt.Println("not enough in stock")
			return -1,-1
		}
		c.moveToFront(c.items[key])
		c.items[key].quantity -= load 
		return key, value-load
	}else{
		//if there's no key in database
		// fmt.Println("there's no key", key, "yet")
		return -1,-1
	}
}

func (c *Cache) put(name string, key int, load int) {
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
	page := c.insert_tail(name, key, load)
	c.size++
	c.items[key] = page
}

func main() {
	cache := newCache(3)
	// cache.put(2, 2)
	// cache.get(2,100)
	// cache.get(1,100)
	// cache.put(1, 1)
	// cache.put(1, 5)
	// cache.get(1, 5)
	// cache.get(2,25)
	// cache.put(8, 8)
	// cache.get(1, 3)
	// cache.get(8, 5)
	// cache.put(2,500)
	// cache.get(2,400)
	// cache.printCache()
	cache.put("toy",1,4)
	cache.printCache()
	cache.put("diamond",2,45)
	cache.printCache()
	cache.put("gucci",23,3247)
	cache.printCache()
	cache.get("diamond",2,30)
	cache.printCache()
	cache.put("car",50,223)
	cache.printCache()

}