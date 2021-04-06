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
	"sync"
	"time"

	// "funcs"
	// "go_work/implem"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

var (
	db    *sql.DB
	mutex sync.Mutex
)

func main() {

	listener, err := net.Listen("tcp", "0.0.0.0:9999")
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

		go handleClientRequest(con)
	}
}

func handleClientRequest(con net.Conn) {
	defer con.Close()

	clientReader := bufio.NewReader(con)
	db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
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
				message = "Insert items"
				number = 1
			} else if clientRequest == "2" {
				log.Println("Remove items")
				message = "Remove items"
				number = 2
			} else if clientRequest == "3" {
				log.Println("Check current stock")
				message = "Check current stock"
				number = 3
			} else if clientRequest == "4" {
				log.Println("Check record for insert")
				message = "Check record for insert"
				number = 4
			} else if clientRequest == "5" {
				log.Println("Check record for remove")
				message = "Check record for remove"
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
			go Going_in(endin, "1", 1, 1)
			// time.Sleep(time.Millisecond)
			<-endin
		} else if number == 2 {
			// db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
			endout := make(chan int)
			go Going_out(endout, "1", 1, 1)
			// time.Sleep(time.Millisecond)
			<-endout
		} else if number == 3 {
			// db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
			endcur := make(chan int)
			go Show_current(endcur, "1")
			// time.Sleep(time.Millisecond)
			<-endcur
		} else if number == 4 {
			// db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
			endrecin := make(chan int)
			go Show_record_in(endrecin, "1")
			// time.Sleep(time.Millisecond)
			<-endrecin
		} else if number == 5 {
			// db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
			endrecout := make(chan int)
			go Show_record_out(endrecout, "1")
			// time.Sleep(time.Millisecond)
			<-endrecout
		}
		// Responding to the client request
		_, err = con.Write([]byte(message + "\n"))
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
	fmt.Println("the id", id, " left in the inventory: ", newQuantity)
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
	start := time.Now()
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
	fmt.Printf("time: %v\n", time.Since(start))

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
	fmt.Println("the id", id, " left in the inventory: ", newQuantity)
	db.Exec("update items set quantity = ? where id = ? ", newQuantity, id)
	c <- 0
}
func Insertingim(n chan string, e chan int, quantity int, id int, name string) {
	product := <-n
	expdate := <-e
	db.Exec("INSERT INTO import(name, quantity, expdate,id,user) VALUES (?, ?, ?, ?, ?)", product, quantity, expdate, id, name)
}

func Going_in(end chan int, name string, id int, quantity int) {
	start := time.Now()
	c := make(chan int)
	q := make(chan int)
	e := make(chan int)
	n := make(chan string)
	if RowExists("SELECT * FROM items WHERE id = ?", id) {
		// fmt.Println(id)
		mutex.Lock()
		go Get_items(q, e, n, id)
		go Plus(q, c, quantity, id)
		<-c // wait for all go routines
		mutex.Unlock()
	} else {
		Insertingitem("New  with id "+strconv.Itoa(id), quantity, 0, id)
	}

	go Insertingim(n, e, quantity, id, name)
	fmt.Printf("time: %v\n", time.Since(start))

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
	return exists
}

func Insertingitem(name string, quantity int, expdate int, id int) {
	db.Exec("INSERT INTO items(name,quantity,expdate,id) VALUES (?,?,?,?)", name, quantity, expdate, id)
}

func Show_current(endrec chan int, name string) {
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
		fmt.Println("name: ", name, "\t quantity: ", quantity, "\t expdate: ", expdate, "\t id: ", id, "\n")
	}
	num, _ := strconv.Atoi(name)
	endrec <- num
}

func Show_record_in(endrec chan int, name string) {
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
		fmt.Println("name: ", name, "\t quantity: ", quantity, "\t expdate: ", expdate, "\t id: ", id, "\t user: ", user, "\n")
	}
	num, _ := strconv.Atoi(name)
	endrec <- num
}

func Show_record_out(endrec chan int, name string) {
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
		fmt.Println("name: ", name, "\t quantity: ", quantity, "\t expdate: ", expdate, "\t id: ", id, "\t user: ", user, "\n")
	}
	num, _ := strconv.Atoi(name)
	endrec <- num
}
