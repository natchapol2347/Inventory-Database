package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db    *sql.DB
	mutex sync.Mutex
)

func checkItem(t chan int, id int) {
	fmt.Println("check item on id:", id)
	row, err := db.Query("select name, quantity, expdate from item where id = " + strconv.Itoa(id))
	if err != nil {
		panic(err)
		//if there is a error on quantity or exp or id then show what is error
	}
	defer row.Close()
	var (
		name     string
		quantity int
		expdate  int
		//create name quatity exp 
	)
	for row.Next() {
		err := row.Scan(&name, &quantity, &expdate)
		if err != nil {
			log.Fatal(err)
			//use log fatal for error handling
		}
		t <- quantity
		fmt.Println("name: ", name, " quantity: ", quantity, " expdate: ", expdate)
	}
}
func export(t chan int, b chan int, exportout int, id int) {
	quantity := <-t // channel from getQuantity
	quantity_left := quantity - exportout
	if quantity_left < 0 {
		b <- 0
		return
	}
	fmt.Println("Item in inventory: ", quantity_left)
	db.Exec("update products set quantity = ? where id = ? ", quantity_left, id)
	b <- 0
}

func impot(t chan int, b chan int, impotin int, id int) {
	quantity := <-t // channel from getQuantity
	quantity_now := quantity + impotin
	fmt.Println("Item in inventory: ", quantity_now)
	db.Exec("update products set quantity = ? where id = ? ", quantity_now, id)
	b <- 0
}

func real_export(stop chan int, id int, exportlist int) {
	db.Exec("update item set quantity = ? where id = ? ", 50, 1)
	b := make(chan int)
	t := make(chan int)
	mutex.Lock()
	go checkItem(t, id)
	go export(t, b, exportlist, id)
	<-b // wait for all go routines
	mutex.Unlock()
}
func real_impot(stop chan int, id int, impotlist int) {
	db.Exec("update item set quantity = ? where id = ? ", 50, 1)
	b := make(chan int)
	t := make(chan int)
	mutex.Lock()
	go checkItem(t, id)
	go impot(t, b, impotlist, id)
	<-b // wait for all go routines
	mutex.Unlock()
}

func main() {
	db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
	defer db.Close()

	n := 50
	stop := make(chan int, n)
	starttime := time.Now()

	// for i := 0; i < n; i++ {
	// 	go real_export(stop, 1, 1)
	// }
	for i := 0; i < n; i++ {
		go real_impot(stop, 2, 1)
	}
	for i := 0; i < n; i++ {
		<-stop
	}
	fmt.Printf("time: %v\n", time.Since(starttime))

	return
}
