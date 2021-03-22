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
	}
	defer row.Close()
	var (
		name     string
		quantity int
		expdate  int
	)
	for row.Next() {
		err := row.Scan(&name, &quantity, &expdate)
		if err != nil {
			log.Fatal(err)
		}
		t <- quantity
		fmt.Println("name: ", name, " quantity: ", quantity, " expdate: ", expdate)
	}
}
func export(t chan int, b chan int, exportout int, id int) {
	quantity := <-t // channel from getQuantity, quantity obtain values from t
	quantity_left := quantity - exportout //current quantity take away exported quantity
	if quantity_left < 0 { //prevent from items being negative
		b <- 0
		return
	}
	fmt.Println("Item in inventory: ", quantity_left)
	db.Exec("update products set quantity = ? where id = ? ", quantity_left, id)
	b <- 0
}

func impot(t chan int, b chan int, impotin int, id int) {
	quantity := <-t // channel from getQuantity, quantity obtain values from t
	quantity_now := quantity + impotin //current quantity add with exported quantity
	fmt.Println("Item in inventory: ", quantity_now)
	db.Exec("update products set quantity = ? where id = ? ", quantity_now, id)
	b <- 0
}

func real_export(stop chan int, id int, exportlist int) {
	//export items in database
	db.Exec("update item set quantity = ? where id = ? ", 50, 1)
	b := make(chan int)
	t := make(chan int)
	//lock to do check first, then export inorder
	mutex.Lock()
	go checkItem(t, id)
	go export(t, b, exportlist, id)
	<-b // wait for all go routines
	mutex.Unlock()
}
func real_impot(stop chan int, id int, impotlist int) {
	//import data into database
	db.Exec("update item set quantity = ? where id = ? ", 50, 1)
	b := make(chan int)
	t := make(chan int)
	//lock to do check first, then import later
	mutex.Lock()
	go checkItem(t, id)
	go impot(t, b, impotlist, id)
	<-b // wait for all go routines
	mutex.Unlock()
}

func main() {
	db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
	//username , password, name of database  
	defer db.Close()

	n := 50
	stop := make(chan int, n)
	starttime := time.Now()

	// for i := 0; i < n; i++ {
	// 	go real_export(stop, 1, 1)
	// }
	for i := 0; i < n; i++ {
		go real_impot(stop, 1, 1) //import concurrently, 50 at the same time
	}
	for i := 0; i < n; i++ {
		<-stop
	}
	fmt.Printf("time: %v\n", time.Since(starttime))

	return
}
