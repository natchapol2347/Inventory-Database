package main

import (

	// "errors"
	// "context"
	"fmt"
	"log"
	"time"

	// "log"
	// "mime/multipart"
	"sync"
	// "time"
	"math/rand"
	"strconv"
)

type Cache struct {
	capacity int
	size int
	items    map[int]*cacheItem
	head     *cacheItem
	tail     *cacheItem
	mu 		sync.Mutex

}

type cacheItem struct {
	name     string
	serial   int
	quantity int
	next     *cacheItem
	prev	 *cacheItem
	mu       sync.Mutex
	last_promoted time.Time

}

func newCache(c int) *Cache {
	return &Cache{
		capacity: c,
		size: 0,
		items:    make(map[int]*cacheItem),
		head:     nil,
		tail: 	  nil,
		mu: 	sync.Mutex{},
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
		last_promoted: time.Time{},

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


func (c *Cache) get(endKey chan int, endQty chan int,name string, key int, load int){
	c.mu.Lock()
	if res, ok := c.items[key]; ok{
		c.mu.Unlock()	
		value := res.quantity
		if(value - load < 0){
			// fmt.Println("not enough in stock")
			
		}
		c.promote(res)
		c.mu.Lock()
		c.items[key].quantity -= load 
		c.mu.Unlock()
		
			endKey <- key
			endQty <- value-load
		
		
		 
	}else{
		c.mu.Unlock()
		//if there's no key
		// fmt.Println("there's no key", key, "yet")
		endKey <- -1
		endQty <- -1
	}
}

func (c *Cache) put(endSig chan int, name string, key int, load int) {
	c.mu.Lock()
	if _, ok := c.items[key]; ok {
		c.mu.Unlock()

		c.mu.Lock()
		c.items[key].quantity += load
		c.mu.Unlock()
		c.promote(c.items[key])
		
		
	}
	c.mu.Unlock()

	if c.size == c.capacity {
		delKey := c.head.serial
		c.pop()
		c.size--
		c.mu.Lock()
		delete(c.items, delKey)
		c.mu.Unlock()
	}
	page := c.insert_tail(name, key, load)
	c.size++
	c.mu.Lock()
	c.items[key] = page
	c.mu.Unlock()
}

func (c *Cache) promote(node *cacheItem) {
	now := time.Now()
	stale := now.Add(time.Minute * -1) // if more than one minute has pass allow for promotion
  
	node.mu.Lock()
	defer node.mu.Unlock()
	if node.last_promoted.Before(stale) {
	  node.last_promoted = now
	  c.mu.Lock()
	  defer c.mu.Unlock()
	  c.moveToFront(node)
	}
	
  }
// func (c *Cache) free(){

// }


func main() {
	
	start := time.Now()
	n := 200
	cache := newCache(n)
	c := make(chan int, n)
	for i:=0; i<n; i++ {
		id := rand.Intn(100)
		name := strconv.Itoa(id)
		go cache.put(c, name,id, 100)

		select{
		case sig := <- c:
			log.Println(sig)
		default:
			log.Println("k")
		}
	

		
	}
	
	// cache.put("toy",1,4)
	// cache.printCache()
	// cache.put("diamond",2,45)
	// cache.printCache()
	// cache.put("gucci",23,3247)
	// cache.printCache()
	// cache.get("diamond",2,30)
	// cache.printCache()
	// cache.put("car",50,223)
	elapsed := time.Since(start)
	log.Println(elapsed)
	// cache.printCache()

}