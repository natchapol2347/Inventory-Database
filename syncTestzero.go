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

	_ "github.com/go-sql-driver/mysql"
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

func (c *Cache) put(end chan int, name string,key int, load int, wg *sync.WaitGroup) {
	//if there's already key in cache just add 
	log.Println("trying1")
	exists := make(chan string)
	full   := make(chan string)
	insert := make(chan string)
	log.Println("trying2")
	// for{
		
		
		go c.keyStateCache(exists, full, insert, key, wg)
		wg.Wait()
		// x:= <-exists
		// y:= <-full
		// z:= <-insert
		log.Println("trying3")
		select{
			case i := <-exists:
				{
					log.Println(i)
					c.mu.Lock()
					c.items[key].quantity += load
					c.mu.Unlock()
					c.mu.Lock()
					c.promote(c.items[key])
					c.mu.Unlock()
					go going_in(nil,nil,end, name, key, load)
					
					
				}
				
			case j := <-full:
				{
					log.Println(j)
					delKey := c.head.serial
					c.pop()
					c.size--
					c.mu.Lock()
					delete(c.items, delKey)
					c.mu.Unlock()
				}
				
			case k := <- insert:
			{
				log.Println(k)
				x := make(chan bool) 
				y := make(chan int)
				go going_in(y, x,end, name, key, load)
				inDB :=<-x
			
				if(!inDB){
					log.Println("not here")
					
				}else{
					log.Println("dsfsdf")
					
					log.Println("heeehhe")
					load := <- y
					
				
					newItem := newItemNode(name, key, load)
					c.lookUp.Lock()
					c.items[key] = newItem
					c.lookUp.Unlock()
					
					c.mu.Lock()
					c.insert_tail(newItem)
					
					c.mu.Unlock()
					
					
					c.size++
					
					

					
					log.Println("les go")
				}
			
			}
			
		
  	//    }
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

func (c *Cache) keyStateCache(existed chan string, full chan string, insert chan string, key int, wg *sync.WaitGroup){
	
	defer wg.Done()
		c.mu.Lock()
		if _, ok := c.items[key]; ok{
			log.Println("yeet1")
			go sendToChannel(existed, "EXISTS")
			// existed <- "EXISTS"
	
		}
		c.mu.Unlock()
		
		if c.size == c.capacity {
			log.Println("yeet2")
			go sendToChannel(full, "FULL")
			// full <- "FULL"
			
		}
		if _, ok := c.items[key]; !ok{
			log.Println("yeet3")
			go sendToChannel(insert, "INSERT")
			// insert <- "INSERT"
			
		}
	
	
	
	
}
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
// 		c.mu.Lock()
// 		c.items[key].quantity -= load 
// 		c.mu.Unlock()
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
// 	c.mu.Lock()
// 	if _, ok := c.items[key]; ok {
// 		log.Println("ficll")
// 		c.mu.Unlock()
// 		c.lookUp.Lock()
// 		c.items[key].quantity += load
// 		c.lookUp.Unlock()
// 		c.promote(c.items[key])
// 		fmt.Println("yoyo")
// 		go going_in(nil,nil,end, name, key, load)
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
	
// 	x := make(chan bool) 
// 	y := make(chan int)

// 	go going_in(y, x,end, name, key, load)
// 	inDB :=<-x
// 	// select{
// 	// 	case inDB := <-x:
// 	// 		{
// 	// 			if(!inDB){
// 	// 				log.Println("not here")
// 	// 				<- end
// 	// 				return 
// 	// 			}
// 	// 		}
// 	// 	default:
// 	// 		{
// 	// 			log.Println("channel not ready")
// 	// 		}
			

// 	// 	}
// 	if(!inDB){
// 		log.Println("not here")
// 		<-end
// 		return

// 	}else{
// 		log.Println("dsfsdf")
// 		load := <- y
// 		log.Println("heeehhe")
		
	
// 		c.mu.Lock()
// 		page := c.insert_tail(name, key, load)
// 		c.mu.Unlock()
// 		c.size++
// 		c.lookUp.RLock()
// 		c.items[key] = page
// 		c.lookUp.RUnlock()


		
// 		log.Println("les go")
		
// 	}
	
// 	// select{
// 	// case load:= <- end:
// 	// 	{
// 	// 		log.Printf("there is %d \n", load)
// 	// 		page := c.insert_tail(name, key, load)
// 	// 		c.size++
// 	// 		c.mu.Lock()
// 	// 		c.items[key] = page
// 	// 		c.mu.Unlock()
// 	// 		log.Println("les go")
// 	// 	}
// 	// default:
// 	// 	log.Printf("waiting to receive")
// 	// }
	
	
	
// }
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

func going_out(sig chan bool, end chan int, name string, id int, quantity int) {
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
		log.Println("burh")
		mutex.Unlock()
	} else {
		sig <- false
		log.Println("burh2")
		return
	}
	go insertingex(n, e, quantity, id, name)
	fmt.Printf("time: %v\n", time.Since(start))

	num, _ := strconv.Atoi(name)
	end <- num
	sig <- true
	fmt.Println("hey?")
	fmt.Println("hello?")
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

func going_in(retQuantity chan int ,sig chan bool, end chan int, name string, id int, quantity int) {
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
		sig <- false
		
	}
	log.Println("wttf")
	go insertingim(n, e, quantity, id, name)
	// fmt.Printf("time: %v\n", time.Since(start))
	log.Println("ee")
	go get_items(q,e,n,id)
	log.Println("aa")

	num := <-q + quantity
	log.Println("oo")

	sig <- true
	log.Println("uu")

	log.Println(num)
	retQuantity <- num

	retName, _:= strconv.Atoi(name)
	end <- retName
	log.Println("haha")

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
	// insertingitem("magazine",122,88,3)
	// insertingitem("sdfasdf",123,23,123)
	result:= make(chan int,100)
	
	// sig:= make(chan *cacheItem, 100)
	cache:= newCache(1000)
	var wg sync.WaitGroup 
	wg.Add(7)
	go cache.put(result, "fruit", 1, 30, &wg)
	go cache.put(result, "1132", 3, 20, &wg)
	go cache.put(result, "223", 1, 20, &wg)
	go cache.put(result, "x", 1, 30, &wg)
	go cache.put(result, "sd", 1, 30, &wg)
	go cache.put(result, "sddsf", 1, 30, &wg)
	go cache.put(result, "s888", 69, 10, &wg)

	

	// cache.get(result, "fu ",6, 1)
	// db.Exec("update items set quantity = ? where id = ? ", 2000, 6)
	// cache.get(result,"444",1,10)
	<- result
	<- result
	<- result
	<- result
	<- result
	<- result
	
	// cache.get(sig, result,"album",23,90)
	// cache.put(sig, "wig",55,50)
	// cache.put(sig, "album",23,100)
	// cache.put(sig, "sfda",123,99)
	cache.printCache()
	fmt.Println(cache.tail)
	// insertingitem("test", 100, 0, 69)
}