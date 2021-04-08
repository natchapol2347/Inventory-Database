package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("====== Start main ======")
	name := "Fake name"

	result := make(chan string)
	go hello(name, result)

	fmt.Println("Finish main")

	fmt.Println(<-result)

	my_res := make(chan string)
	go checkin("123", 2, my_res)

	fmt.Println(<-my_res)
}

func hello(name string, result chan<- string) {
	output := "Hello " + name
	fmt.Printf("In function = %s\n", output)
	result <- output
}

func checkin(itmNo string, q int16, result chan<- string) {
	start2 := time.Now() //start what time
	fmt.Print(start2)    //display time
	itemCode := itmNo
	qty := q

	// print display on screen
	fmt.Printf("Item added %s qty=%d\n", itemCode, qty)
	start2 = time.Now() //start what time
	fmt.Print(start2)   //display time

	my_res := "OK : MyRes"

	result <- my_res
} //. End checkin
