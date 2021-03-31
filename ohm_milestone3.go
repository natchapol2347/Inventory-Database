package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
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

func going_out(end chan int, name string, id int, quantity int) {
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
		return
	}
	go insertingex(n, e, quantity, id, name)
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
	fmt.Println("the items left in stock: ", newQuantity)
	db.Exec("update items set quantity = ? where id = ? ", newQuantity, id)
	c <- 0
}
func insertingim(n chan string, e chan int, quantity int, id int, name string) {
	product := <-n
	expdate := <-e
	db.Exec("INSERT INTO import(name, quantity, expdate,id,user) VALUES (?, ?, ?, ?, ?)", product, quantity, expdate, id, name)
}

func going_in(end chan int, name string, id int, quantity int) {
	start := time.Now()
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
		insertingitem("New items", quantity, 0, id)
	}

	go insertingim(n, e, quantity, id, name)
	fmt.Printf("time: %v\n", time.Since(start))

	num, _ := strconv.Atoi(name)
	end <- num
	return
}

func rowExists(query string, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		glog.Fatalf("error checking if row exists '%s' %v", args, err)
	}
	return exists
}

func insertingitem(name string, quantity int, expdate int, id int) {
	db.Exec("INSERT INTO items(name,quantity,expdate,id) VALUES (?,?,?,?)", name, quantity, expdate, id)
}

func show_current() {
	rows, err := db.Query("SELECT * FROM items")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var name string
		var quantity int
		var expdate int
		var id int
		err = rows.Scan(&name, &quantity, &expdate, &id)
		if err != nil {
			panic(err)
		}
		fmt.Println("name: ", name, " quantity: ", quantity, " expdate: ", expdate, "id: ", id,"\n")
	}
}

func show_record_in() {
	rows, err := db.Query("SELECT * FROM import")
	if err != nil {
		panic(err)
	}
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
		fmt.Println("name: ", name, " quantity: ", quantity, " expdate: ", expdate, "id: ", id,"user: ", user,"\n")
	}
}

func show_record_out() {
	rows, err := db.Query("SELECT * FROM export")
	if err != nil {
		panic(err)
	}
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
		fmt.Println("name: ", name, " quantity: ", quantity, " expdate: ", expdate, "id: ", id,"user: ", user,"\n")
	}
}

func main() {
	db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
	//defer db.Close()
	// insertingitem("TV",1000,0,3)
	n := 50
	endin := make(chan int, n)
	endout := make(chan int, n)
	start2 := time.Now()

	// for i := 0; i < n; i++ {
	// 	go show_record_in()
	// 	go show_record_out()
	// 	go show_current()
	// }
	for i := 0; i < n; i++ {
		go going_out(endout, strconv.Itoa(i), 3, 1)
	}
	for i := 0; i < n; i++ {
		go going_in(endin, strconv.Itoa(i), 1, 1)
	}
	for i := 0; i < n; i++ {
		<-endin
	}
	for i := 0; i < n; i++ {
		<-endout
	}
	fmt.Printf("Total time: %v\n", time.Since(start2))
	
	return
}
