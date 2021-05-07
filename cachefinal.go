package main

import (

	// "errors"
	"fmt"
	// "math/big"
	// "mime/multipart"
	"database/sql"
	"log"
	"strconv"
	"sync"
	"time"

	_"github.com/go-sql-driver/mysql"
	// "github.com/golang/glog"
)

//Cache codes
type Cache struct {
	capacity int
	size int
	items    map[int]*cacheItem
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
		items:    make(map[int]*cacheItem),
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
func (c *Cache) insert_tail(newItem *cacheItem) {
	//make new item from argument
	
	if(c.tail == nil && c.head == nil){
		c.tail = newItem
		c.head = newItem

	}else{
		newItem.next = c.tail
		c.tail.prev = newItem
		c.tail = newItem
	}
	
	

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

func (c *Cache) put(end chan int, name string,key int, load int) {
	//if there's already key in cache just add 
	exists := make(chan string)
	full   := make(chan string)
	insert := make(chan string)
		go c.keyStateCache(exists, full, insert, key)
		select{
			case <-exists:
				{
					c.mu.Lock()
					c.items[key].quantity += load
					c.mu.Unlock()
					c.mu.Lock()
					c.promote(c.items[key])
					c.mu.Unlock()
					go going_in(nil,end, name, key, load)
				}
			case <-full:
				{
					delKey := c.head.serial
					c.pop()
					c.size--
					c.mu.Lock()
					delete(c.items, delKey)
					c.mu.Unlock()
					go sendToChannel(insert,"INSERT")
				}	
			case <- insert:
			{
				y := make(chan int)
				go going_in(y,end, name, key, load)
				load := <- y
				newItem := newItemNode(name, key, load)
				c.lookUp.Lock()
				c.items[key] = newItem
				c.lookUp.Unlock()
				c.mu.Lock()
				c.insert_tail(newItem)
				c.mu.Unlock()
				c.size++
			}
	}
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

func (c *Cache) keyStateCache(existed chan string, full chan string, insert chan string, key int){
	
	
		c.mu.Lock()
		if _, ok := c.items[key]; ok{
			
			go sendToChannel(existed, "EXISTS")
			// existed <- "EXISTS"
	
		}
		c.mu.Unlock()
		
		if c.size == c.capacity {
			
			go sendToChannel(full, "FULL")
			// full <- "FULL"
			
		}
		if _, ok := c.items[key]; !ok{
			
			go sendToChannel(insert, "INSERT")
			// insert <- "INSERT"
			
		}
	
	
	
	
}
func (c *Cache) get(end chan int, name string, key int, load int) {
	exists := make(chan string)
	full   := make(chan string)
	insert := make(chan string)
	go c.keyStateCache(exists, full, insert, key)
	select{
		case <- exists:
			{
				c.lookUp.RLock()			
				if(c.items[key].quantity-load<0){
					log.Println("cache: quantity less than request")
					end <- 0
				}
				c.lookUp.RUnlock()
				
				c.mu.Lock()
				c.items[key].quantity -= load 
				c.mu.Unlock()
				c.lookUp.Lock()
				c.promote(c.items[key])
				c.lookUp.Unlock()
				go going_out(nil,nil, end, name, key, load)
			}
		case <-full:
			{
				delKey := c.head.serial
				c.pop()
				c.size--
				c.mu.Lock()
				delete(c.items, delKey)
				c.mu.Unlock()	
				go sendToChannel(insert, "INSERT")
			}
			
		case <- insert:
		{
			x := make(chan bool) 
			y := make(chan int)
			go going_out(y, x,end, name, key, load)
			inDB :=<-x
			if(!inDB){
				log.Println("cache: item not in db")
				<-y
				end <- 0

			}else{//cache miss
				
				load := <- y
				newItem := newItemNode(name, key, load)
				c.lookUp.Lock()
				c.items[key] = newItem
				c.lookUp.Unlock()
				c.mu.Lock()
				c.insert_tail(newItem)
				c.mu.Unlock()
				c.size++
			}
		}	
	}
}
func sendToChannel(ch chan string, sig string){
	ch <- sig
}

func dumpChannel(ch chan int){
	<- ch
}

func (c *Cache) promote(node *cacheItem) {
	now := time.Now()
	stale := now.Add(time.Minute * -1) // if more than one minute has pass allow for promotion
  
	c.mu.Lock()
	defer c.mu.Unlock()
	if node.last_promoted.Before(stale) {
	  node.last_promoted = now
	  c.mu.Lock()
	  defer c.mu.Unlock()
	  c.moveToFront(node)
	}
	
  }

//Database code
var (
	db    *sql.DB
	mutex sync.Mutex
)

func get_items(q chan int, e chan int, n chan string, id int) {

	row, err := db.Query("SELECT name, quantity, expdate FROM items WHERE id = " + strconv.Itoa(id))
	if err != nil {
		panic(err)
	}
	defer row.Close()
	for row.Next() {
		var (
			name     string
			quantity int
			expdate  int
		)
		row.Scan(&name, &quantity, &expdate)
		q <- quantity
		n <- name
		e <- expdate
		
	}
}
func decrement(q chan int, c chan int, quantity int, id int) {
	quantityy := <-q // channel from get_items
	newQuantity := quantityy - quantity
	if newQuantity < 0 {
		c <- 0
		return
	}
	fmt.Println("the items left in stock: ", newQuantity)
	db.Exec("update items set quantity = ? where id = ? ", newQuantity, id)
	log.Println(newQuantity)
	c <- 0
}

func insertingex(n chan string, e chan int, quantity int, id int, name string) {
	product := <-n
	expdate := <-e
	db.Exec("INSERT INTO export(name, quantity, expdate,id,user) VALUES (?, ?, ?, ?, ?)", product, quantity, expdate, id, name)
}

func going_out(retQuantity chan int, sig chan bool, end chan int, name string, id int, quantity int) {
	
	c := make(chan int)
	q := make(chan int)
	e := make(chan int)
	n := make(chan string)
	if rowExists("SELECT * FROM items WHERE id = ?", id) {
		mutex.Lock()
		go get_items(q, e, n, id)
		go decrement(q, c, quantity, id)
		<-c // wait for all go routines
		mutex.Unlock()
	} else {
		sig <- false
		retQuantity <- 0
		return
	}
	go insertingex(n, e, quantity, id, name)
	go get_items(q, e, n, id)
	num:= (<-q -quantity)
	sig <-true
	retQuantity <- num

	retName, _:= strconv.Atoi(name)
	end <- retName
}

func increment(q chan int, c chan int, quantity int, id int) {
	quantityy := <-q // channel from get_items
	newQuantity := quantityy + quantity
	if newQuantity < 0 {
		c <- 0
		return
	}
	// fmt.Println("the items left in stock: ", newQuantity)
	db.Exec("update items set quantity = ? where id = ? ", newQuantity, id)
	c <- 0
}
func insertingim(n chan string, e chan int, quantity int, id int, name string) {
	product := <-n
	expdate := <-e
	db.Exec("INSERT INTO import(name, quantity, expdate,id,user) VALUES (?, ?, ?, ?, ?)", product, quantity, expdate, id, name)
}

func going_in(retQuantity chan int, end chan int, name string, id int, quantity int) {
	c := make(chan int)
	q := make(chan int)
	e := make(chan int)
	n := make(chan string)
	
	if rowExists("SELECT * FROM items WHERE id = ?", id) {
		mutex.Lock()
		go get_items(q, e, n, id)
		go increment(q, c, quantity, id)
		<-c // wait for all go routines
		mutex.Unlock()
	} else {
		insertingitem("New  with id "+strconv.Itoa(id), quantity, 0, id)
		
		
	}
	go insertingim(n, e, quantity, id, name)
	go get_items(q,e,n,id)
	num := <-q + quantity


	retQuantity <- num

	retName, _:= strconv.Atoi(name)
	end <- retName


}

func rowExists(query string, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal("error checking if row exists '%s' %v", args, err)
	}
	return exists
}

func insertingitem(name string, quantity int, expdate int, id int) {
	
	db.Exec("INSERT INTO items(name,quantity,expdate,id) VALUES (?,?,?,?)", name, quantity, expdate, id)
}

func show_current(endrec chan int, name string) {
	rows, err := db.Query("SELECT * FROM items")
	if err != nil {
		panic(err)
	}
	fmt.Println("current items: \n")
	for rows.Next() {
		var name string
		var quantity int
		var expdate int
		var id int
		err = rows.Scan(&name, &quantity, &expdate, &id)
		if err != nil {
			panic(err)
		}
		fmt.Println("name: ", name, "\t quantity: ", quantity, "\t expdate: ", expdate, "\t id: ", id,"\n")
	}
	num, _ := strconv.Atoi(name)
	endrec<-num
}

func show_record_in(endrec chan int, name string) {
	rows, err := db.Query("SELECT * FROM import")
	if err != nil {
		panic(err)
	}
	fmt.Println("record in: \n")
	for rows.Next() {
		var name string
		var quantity int
		var expdate int
		var id int
		var user string
		err = rows.Scan(&name, &quantity, &expdate, &id, &user)
		if err != nil {
			panic(err)
		}
		fmt.Println("name: ", name, "\t quantity: ", quantity, "\t expdate: ", expdate, "\t id: ", id,"\t user: ", user,"\n")
	}
	num, _ := strconv.Atoi(name)
	endrec<-num
}

func show_record_out(endrec chan int, name string) {
	rows, err := db.Query("SELECT * FROM export")
	if err != nil {
		panic(err)
	}
	fmt.Println("record out: \n")
	for rows.Next() {
		var name string
		var quantity int
		var expdate int
		var id int
		var user string
		err = rows.Scan(&name, &quantity, &expdate, &id, &user)
		if err != nil {
			panic(err)
		}
		fmt.Println("name: ", name, "\t quantity: ", quantity, "\t expdate: ", expdate, "\t id: ", id,"\t user: ", user,"\n")
	}
	num, _ := strconv.Atoi(name)
	endrec<-num
}

func main(){
	db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
	defer db.Close()
	result:= make(chan int,100)	
	cache:= newCache(10)
	go cache.put(result, "fruit", 1, 30)
	go cache.put(result, "car", 3, 20)
	go cache.put(result, "fruit", 1, 20)
	go cache.put(result, "fruit", 1, 30)
	go cache.put(result, "fruit", 1, 30)
	go cache.put(result, "camera", 22, 30)
	go cache.get(result, "cd", 69, 10)
	go cache.get(result, "camera", 22, 100)
	<- result
	<- result
	<- result
	<- result
	<- result
	<- result
	<- result
	cache.printCache()
}