package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"math/rand"
	"sync"
	// "time"

	// "funcs"
	// "go_work/implem"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

var (
	db    *sql.DB
	mutex sync.Mutex
	count int = 0
)

func main() {

	listener, err := net.Listen("tcp", "127.0.0.3:8888")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		count++
		fmt.Println("clients ", count)
		go handleClientRequest(con)
	}
}

func handleClientRequest(con net.Conn) {
	defer con.Close()

	clientReader := bufio.NewReader(con)
	db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
	db.SetMaxOpenConns(100000)
	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')
		message := "Please provide numbers 1-5"
		number := 0
		switch err {
		case nil:

			clientRequest := strings.TrimSpace(clientRequest)
			if clientRequest == "QUIT" {
				log.Println("client requested server to close the connection so closing")
				return
			} else if clientRequest == "1" {
				log.Println("Insert items")
				// message = "Insert items"
				number = 1
			} else if clientRequest == "2" {
				log.Println("Remove items")
				// message = "Remove items"
				number = 2
			} else if clientRequest == "3" {
				log.Println("Check current stock")
				// message = "Check current stock"
				number = 3
			} else if clientRequest == "4" {
				log.Println("Check record for insert")
				// message = "Check record for insert"
				number = 4
			} else if clientRequest == "5" {
				log.Println("Check record for remove")
				// message = "Check record for remove"
				number = 5
			} else {
				log.Println("Please provide numbers 1-5")
			}
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}

		if number == 1 {
			// db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
			endin := make(chan int)
			go Going_in(endin, strconv.Itoa(rand.Intn(10) + 1), rand.Intn(10) + 1, rand.Intn(10) + 1)
			// time.Sleep(time.Millisecond)
			message = "The item has been added."
			<-endin
		} else if number == 2 {
			// db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
			endout := make(chan int)
			go Going_out(endout, strconv.Itoa(rand.Intn(10) + 1), rand.Intn(10) + 1, rand.Intn(10) + 1)
			// time.Sleep(time.Millisecond)
			message = "The item has been removed."
			<-endout
		} else if number == 3 {
			// db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
			endcur := make(chan int)
			me := make(chan string)
			go Show_current(endcur, strconv.Itoa(rand.Intn(10) + 1),me)
			// time.Sleep(time.Millisecond)
			message = <-me
			<-endcur
		} else if number == 4 {
			// db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
			endrecin := make(chan int)
			me := make(chan string)
			go Show_record_in(endrecin, strconv.Itoa(rand.Intn(10) + 1),me)
			// time.Sleep(time.Millisecond)
			message = <-me
			<-endrecin
		} else if number == 5 {
			// db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
			endrecout := make(chan int)
			me := make(chan string)
			go Show_record_out(endrecout, strconv.Itoa(rand.Intn(10) + 1),me)
			// time.Sleep(time.Millisecond)
			message = <-me
			<-endrecout
		}
		// Responding to the client request
		_, err = con.Write([]byte(message+"\n"))
		if err != nil {
			log.Printf("failed to respond to client: %v\n", err)
		}

	}
}

func Get_items(q chan int, e chan int, n chan string, id int) {

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
		// fmt.Println("name: ", name, " quantity: ", quantity, " expdate: ", expdate)
	}
}
func Minus(q chan int, c chan int, quantity int, id int) {
	quantityy := <-q // channel from get_items
	newQuantity := quantityy - quantity
	if newQuantity < 0 {
		c <- 0
		return
	}
	// fmt.Println("the id", id, " left in the inventory: ", newQuantity)
	db.Exec("update items set quantity = ? where id = ? ", newQuantity, id)
	c <- 0
}

func Insertingex(n chan string, e chan int, quantity int, id int, name string) {
	product := <-n
	expdate := <-e
	db.Exec("INSERT INTO export(name, quantity, expdate,id,user) VALUES (?, ?, ?, ?, ?)", product, quantity, expdate, id, name)
}

func Going_out(end chan int, name string, id int, quantity int) {
	// fmt.Printf("start\n")
	// start := time.Now()
	c := make(chan int)
	q := make(chan int)
	e := make(chan int)
	n := make(chan string)
	if RowExists("SELECT * FROM items WHERE id = ?", id) {
		mutex.Lock()
		go Get_items(q, e, n, id)
		go Minus(q, c, quantity, id)
		<-c // wait for all go routines
		mutex.Unlock()
	} else {
		return
	}
	go Insertingex(n, e, quantity, id, name)
	// fmt.Printf("time: %v\n", time.Since(start))

	num, _ := strconv.Atoi(name)
	end <- num
	return
}

func Plus(q chan int, c chan int, quantity int, id int) {
	quantityy := <-q // channel from get_items
	newQuantity := quantityy + quantity
	if newQuantity < 0 {
		c <- 0
		return
	}
	// fmt.Println("the id", id, " left in the inventory: ", newQuantity)
	db.Exec("update items set quantity = ? where id = ? ", newQuantity, id)
	c <- 0
}
func Insertingim(n chan string, e chan int, quantity int, id int, name string) {
	product := <-n
	expdate := <-e
	db.Exec("INSERT INTO import(name, quantity, expdate,id,user) VALUES (?, ?, ?, ?, ?)", product, quantity, expdate, id, name)
}

func insertingim_first(name string, expdate int, quantity int, id int, nameuse string){
	db.Exec("INSERT INTO import(name, quantity, expdate,id,user) VALUES (?, ?, ?, ?, ?)", name, quantity, expdate, id, nameuse)
}

func Going_in(end chan int, name string, id int, quantity int) {
	// start := time.Now()
	// d := make(chan int)
	c := make(chan int)
	q := make(chan int)
	e := make(chan int)
	n := make(chan string)
	mutex.Lock()
	if RowExists("SELECT * FROM items WHERE id = ?", id) {
		fmt.Println("yeeha")
		// mutex.Lock()
		go Get_items(q, e, n, id)
		go Plus(q, c, quantity, id)
		// go Insertingim(n, e, quantity, id, name,c)
		<-c // wait for all go routines
		// mutex.Unlock()
		// go Insertingim(n, e, quantity, id, name)
	} else {
		
		// mutex.Lock()
		fmt.Println("hello")
		go Insertingitem("New  with id "+strconv.Itoa(id), quantity, 0, id,c)
		// go Get_items(q, e, n, id)
		// go Plus(q, c, quantity, id)
		go insertingim_first("New  with id "+strconv.Itoa(id), 0, quantity, id, name)
		<-c // wait for all go routines
		// mutex.Unlock()
	}
	mutex.Unlock()
	go Insertingim(n, e, quantity, id, name)
	// go Insertingim(n, e, quantity, id, name)
	// fmt.Printf("time: %v\n", time.Since(start))
	
	num, _ := strconv.Atoi(name)
	end <- num
	return
}

func RowExists(query string, args ...interface{}) bool {
	var exists bool
	// fmt.Println(args)
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		glog.Fatalf("error checking if row exists '%s' %v", args, err)
	}
	// d<-0
	// <-d
	return exists
}

func Insertingitem(name string, quantity int, expdate int, id int, c chan int) {
	// mutex.Lock()
	// mutex.Unlock()
	// namee:=<-n
	db.Exec("INSERT INTO items(name,quantity,expdate,id) VALUES (?,?,?,?)", name, quantity, expdate, id)
	// go Insertingim(n, e, quantity, id, namee)
	c<-0
	
}

func Show_current(endrec chan int, name string, me chan string) {
	rows, err := db.Query("SELECT * FROM items")
	if err != nil {
		panic(err)
	}
	whole := "current items: \n"
	for rows.Next() {
		var namee string
		var quantity int
		var expdate int
		var id int
		err = rows.Scan(&name, &quantity, &expdate, &id)
		if err != nil {
			panic(err)
		}
		line := "name: "+ namee+ " quantity: "+ strconv.Itoa(quantity)+ " expdate: "+ strconv.Itoa(expdate)+ " id: "+strconv.Itoa(id)+"\n"
		whole = whole+line
	}
	whole = whole+"."
	fmt.Println(whole)
	me<-whole
	num, _ := strconv.Atoi(name)
	endrec <- num
}

func Show_record_in(endrec chan int, name string,me chan string) {
	rows, err := db.Query("SELECT * FROM import")
	if err != nil {
		panic(err)
	}
	whole := "record for import: \n"
	for rows.Next() {
		var namee string
		var quantity int
		var expdate int
		var id int
		var user string
		err = rows.Scan(&namee, &quantity, &expdate, &id, &user)
		if err != nil {
			panic(err)
		}
		line := "name: "+ namee+ " quantity: "+ strconv.Itoa(quantity)+ " expdate: "+ strconv.Itoa(expdate)+ " id: "+strconv.Itoa(id)+" user: "+user+"\n"
		whole = whole+line
	}
	whole = whole+"."
	fmt.Println(whole)
	me<-whole
	num, _ := strconv.Atoi(name)
	endrec <- num
}

func Show_record_out(endrec chan int, name string,me chan string) {
	rows, err := db.Query("SELECT * FROM export")
	if err != nil {
		panic(err)
	}
	whole := "record for export: \n"
	for rows.Next() {
		var namee string
		var quantity int
		var expdate int
		var id int
		var user string
		err = rows.Scan(&namee, &quantity, &expdate, &id, &user)
		if err != nil {
			panic(err)
		}
		line := "name: "+ namee+ " quantity: "+ strconv.Itoa(quantity)+ " expdate: "+ strconv.Itoa(expdate)+ " id: "+strconv.Itoa(id)+" user: "+user+"\n"
		whole = whole+line
	}
	whole = whole+"."
	fmt.Println(whole)
	me<-whole
	num, _ := strconv.Atoi(name)
	endrec <- num
}

// //Cache codes
// type Cache struct {
// 	capacity int
// 	size int
// 	items    map[int]*cacheItem
// 	mu       sync.Mutex
// 	head     *cacheItem
// 	tail     *cacheItem
// }

// type cacheItem struct {
// 	name     string
// 	serial   int
// 	quantity int
// 	next     *cacheItem
// 	prev	 *cacheItem
// 	mu       sync.Mutex
// 	last_promoted time.Time
// }

// func newCache(c int) *Cache {
// 	return &Cache{
// 		capacity: c,
// 		size: 0,
// 		items:    make(map[int]*cacheItem),
// 		mu:       sync.Mutex{},
// 		head:     nil,
// 		tail: 	  nil,
// 	}
// }

// func newItemNode(in_name string, key int, value int) *cacheItem{
// 	return &cacheItem{
// 		name: in_name,
// 		serial: key,
// 		quantity: value,
// 		next: nil,
// 		prev: nil,
// 		mu :     sync.Mutex{},
// 		last_promoted: time.Time{},
// 	}
// }
// func (c *Cache) insert_tail(name string,key int, value int) *cacheItem{
// 	//make new item from argument
// 	newItem := newItemNode(name, key, value)
// 	if(c.tail == nil && c.head == nil){
// 		c.tail = newItem
// 		c.head = newItem

// 	}else{
// 		newItem.next = c.tail
// 		c.tail.prev = newItem
// 		c.tail = newItem
// 	}

// 	return newItem

// }

// func (c *Cache) moveToFront(node *cacheItem){
// 	if node == c.tail{
// 		return 
// 	}else if node == c.head{
// 		c.head = c.head.prev
// 		//last node's next must point to nil
// 		c.head.next = nil 
// 	}else{
// 		node.prev.next = node.next
// 		node.next.prev = node.prev
// 	}

// 	node.next = c.tail
// 	c.tail.prev = node
// 	c.tail = node
// }

// func (c *Cache) pop(){
// 	if c.head == nil && c.tail == nil{
// 		return
// 	}else if c.head == c.tail{
// 		c.head, c.tail = nil, nil
// 	}else{
// 		c.head = c.head.prev
// 		c.head.next = nil
// 	}
// }
// func (c *Cache) printCache(){
// 	current := c.tail
// 	var i int;
// 	for i=0;i<=c.size;i++{
// 		if(current != nil){		
// 			fmt.Printf("|name:%s|id:%d|,quantity:%d|size:%d| ->", current.name, current.serial, current.quantity,c.size)
// 			current = current.next
// 		}
// 	}
// 	fmt.Println("\n")
// }


// func (c *Cache) get(end chan int, name string, key int, load int) {
// 	c.mu.Lock()
// 	if res, ok := c.items[key]; ok{
// 		c.mu.Unlock()
// 		value := res.quantity
// 		if(value - load < 0){
// 			// fmt.Println("not enough in stock")
// 			<- end
// 			return 
// 		}
// 		c.promote(res)
// 		res.mu.Lock()
// 		c.items[key].quantity -= load 
// 		res.mu.Unlock()
// 		// c.mu.Lock()
// 		// result <- c.items[key]
// 		// c.mu.Unlock()
// 		go going_out(nil, end, name, key, load)
// 		<- end
		
// 		return 
// 	}else{
// 		//if there's no key
// 		// fmt.Println("there's no key", key, "yet")
// 		c.mu.Unlock()
// 		x := make(chan bool) 
// 		defer close(x)
		
// 		go going_out(x, end, name, key, load)
		
// 		// select{
// 		// case inDB := <-x:
// 		// 	{
// 		// 		if(!inDB){
// 		// 			log.Println("not here")
// 		// 			<- end
// 		// 			return 
// 		// 		}
// 		// 	}
// 		// default:
// 		// 	{
// 		// 		log.Println("channel not ready")
// 		// 	}
			

// 		// }
// 		if(!<-x){
// 			log.Println("not here")
// 			<-end
// 			return

// 		}else{
// 			log.Println("channel not ready")
// 		}

// 		log.Println("kksdf")
// 		c.put(end, name, key, load)
// 		// c.mu.Lock()
// 		// result <- c.items[key]
// 		// c.mu.Unlock()
// 		fmt.Println("kksdf")
// 		<- end
// 		return 
		
// 	}
	
// }

// func (c *Cache) put(end chan int, name string,key int, load int) {
// 	//if there's already key in cache just add 
// 	c.mu.Lock()
// 	if res, ok := c.items[key]; ok {
// 		c.mu.Unlock()
// 		res.mu.Lock()
// 		res.quantity += load
// 		res.mu.Unlock()
// 		c.promote(c.items[key])
// 		fmt.Println("yoyo")
// 		go going_in(end, name, key, load)
// 		<- end
// 		return
// 	}
// 	c.mu.Unlock()
// 	//if cache is full delete last recent node
// 	if c.size == c.capacity {
// 		delKey := c.head.serial
// 		c.pop()
// 		c.size--
// 		c.mu.Lock()
// 		delete(c.items, delKey)
// 		c.mu.Unlock()
// 	}
// 	//insert new node
// 	page := c.insert_tail(name, key, load)
// 	c.size++
// 	c.mu.Lock()
// 	c.items[key] = page
// 	c.mu.Unlock()
// 	go going_in(end, name, key, load)
// 	<- end
// 	return

	

// }

// func (c *Cache) promote(node *cacheItem) {
// 	now := time.Now()
// 	stale := now.Add(time.Minute * -1) // if more than one minute has pass allow for promotion
  
// 	node.mu.Lock()
// 	defer node.mu.Unlock()
// 	if node.last_promoted.Before(stale) {
// 	  node.last_promoted = now
// 	  c.mu.Lock()
// 	  defer c.mu.Unlock()
// 	  c.moveToFront(node)
// 	}
	
//   } 
