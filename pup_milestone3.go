package main

import (

	// "errors"
	"fmt"
	// "mime/multipart"
	"sync"
	"database/sql"
	"strconv"
	"time"
	"log"
	_"github.com/go-sql-driver/mysql"
	// "github.com/golang/glog"

)

//Cache codes
type Cache struct {
	capacity int
	size int
	items    map[int]*cacheItem
	mu       sync.Mutex
	head     *cacheItem
	tail     *cacheItem
}

type cacheItem struct {
	name     string
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

func newItemNode(in_name string, key int, value int) *cacheItem{
	return &cacheItem{
		name: in_name,
		serial: key,
		quantity: value,
		next: nil,
		prev: nil,
	}
}
func (c *Cache) insert_tail(name string,key int, value int) *cacheItem{
	//make new item from argument
	newItem := newItemNode(name, key, value)
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


func (c *Cache) get(end chan int, name string, key int, load int) (string, int,int){
	if _, ok := c.items[key]; ok{
		value := c.items[key].quantity
		if(value - load < 0){
			// fmt.Println("not enough in stock")
			return "",-1,-1
		}
		c.moveToFront(c.items[key])
		c.items[key].quantity -= load 
		go going_out(nil,end, name, key, load)
		return name, key, value-load
	}else{
		//if there's no key
		// fmt.Println("there's no key", key, "yet")
		x := make(chan bool) 
		defer close(x)
		go going_out(x, end, name, key, load)
		var ifIn bool = <-x 
		if(!ifIn){
			fmt.Println("not here!")
			return "",-1,-1
		}
		c.put(end, name, key, load)
		return name,key, c.items[key].quantity
	}
}

func (c *Cache) put(end chan int, name string,key int, load int) {
	//if there's already key in cache just add 
	if _, ok := c.items[key]; ok {
		c.items[key].quantity += load
		c.moveToFront(c.items[key])
		go going_in(end, name, key, load)
		return
	}
	//if cache is full delete last recent node
	if c.size == c.capacity {
		delKey := c.head.serial
		c.pop()
		c.size--
		delete(c.items, delKey)
	}
	//insert new node
	page := c.insert_tail(name, key, load)
	c.size++
	c.items[key] = page
	go going_in(end, name, key, load)
	return

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
	c <- 0
}

func insertingex(n chan string, e chan int, quantity int, id int, name string) {
	product := <-n
	expdate := <-e
	db.Exec("INSERT INTO export(name, quantity, expdate,id,user) VALUES (?, ?, ?, ?, ?)", product, quantity, expdate, id, name)
}

func going_out(cacheEnd chan bool, end chan int, name string, id int, quantity int) {
	// fmt.Printf("start\n")
	start := time.Now()
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
		cacheEnd <- false
		return
	}
	go insertingex(n, e, quantity, id, name)
	fmt.Printf("time: %v\n", time.Since(start))

	num, _ := strconv.Atoi(name)
	end <- num
	cacheEnd <-true
	return
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

func going_in(end chan int, name string, id int, quantity int) {
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
		fmt.Println("adding new item")
		insertingitem("New  with id "+strconv.Itoa(id), quantity, 0, id)
		
		
	}
	go insertingim(n, e, quantity, id, name)
	// fmt.Printf("time: %v\n", time.Since(start))

	num, _ := strconv.Atoi(name)
	end <- num
	
	
	return 
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
	// insertingitem("fruit",100,29,3)
	c:= make(chan int,100)

	cache:= newCache(5)

	cache.put(c,"fruit", 1, 30)
	cache.get(c,"wig",55, 40)
	cache.get(c,"fruit",1,10)
	cache.get(c,"album",23,90)
	cache.put(c,"wig",55,50)
	cache.put(c,"album",23,100)
	<- c
	cache.printCache()
}