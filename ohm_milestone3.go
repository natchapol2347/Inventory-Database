package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

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
		fmt.Println("name: ", name, " quantity: ", quantity, " expdate: ", expdate)
	}
}
func decrement(q chan int, c chan int, orderQuantity int, id int) {
	quantity := <-q // channel from get_items
	newQuantity := quantity - orderQuantity
	if newQuantity < 0 {
		c <- 0
		return
	}
	fmt.Println("Product left in stock: ", newQuantity)
	db.Exec("update items set quantity = ? where id = ? ", newQuantity, id)
	c <- 0
}

func insertingex(n chan string, e chan int, quantity int, id int, name string) {
	product := <-n
	expdate := <-e
	db.Exec("INSERT INTO export(name, quantity, expdate,id,user) VALUES (?, ?, ?, ?, ?)", product, quantity, expdate, id, name)
}

func going_out(end chan int, name string, productId int, orderQuantity int) {
	// fmt.Printf("start\n")
	start := time.Now()
	c := make(chan int)
	q := make(chan int)
	e := make(chan int)
	n := make(chan string)
	mutex.Lock()
	go get_items(q, e, n, productId)
	go decrement(q, c, orderQuantity, productId)
	<-c // wait for all go routines
	mutex.Unlock()
	go insertingex(n, e, orderQuantity, productId, name)
	fmt.Printf("time: %v\n", time.Since(start))

	num, _ := strconv.Atoi(name)
	end <- num
	return
}

func increment(q chan int, c chan int, quantity int, id int) {
	quantityy := <-q // channel from get_items
	newQuantity := quantityy + quantity
	if newQuantity < 0 {
		c <- 0
		return
	}
	fmt.Println("Product left in stock: ", newQuantity)
	db.Exec("update items set quantity = ? where id = ? ", newQuantity, id)
	c <- 0
}
func insertingim(n chan string, e chan int, quantity int, id int, name string) {
	product := <-n
	expdate := <-e
	db.Exec("INSERT INTO import(name, quantity, expdate,id,user) VALUES (?, ?, ?, ?, ?)", product, quantity, expdate, id, name)
}

func going_in(end chan int, name string, productId int, orderQuantity int) {
	start := time.Now()
	c := make(chan int)
	q := make(chan int)
	e := make(chan int)
	n := make(chan string)
	mutex.Lock()
	go get_items(q, e, n, productId)
	go increment(q, c, orderQuantity, productId)
	<-c // wait for all go routines
	mutex.Unlock()
	go insertingim(n, e, orderQuantity, productId, name)
	fmt.Printf("time: %v\n", time.Since(start))

	num, _ := strconv.Atoi(name)
	end <- num
	return
}

func insertingitem(name string, quantity int, expdate int, id int) {
	db.Exec("INSERT INTO items(name,quantity,expdate,id) VALUES (?,?,?,?)", name, quantity, expdate, id)
}

func main() {
	db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
	//defer db.Close()
	// insertingitem("TV",1000,0,3)
	n := 100
	end := make(chan int, n)
	start2 := time.Now()

	for i := 0; i < n; i++ {
		go going_out(end, strconv.Itoa(i), 1, 1)
		go going_in(end, strconv.Itoa(i), 1, 1)
	}
	for i := 0; i < n; i++ {
		<-end
	}
	fmt.Printf("Total time: %v\n", time.Since(start2))

	return
}
