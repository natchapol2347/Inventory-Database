package main

import (
	"fmt"
	"time"
)

func checkin(itmNo string, q int16, result chan<- string) {
	start := time.Now() //start what time
	fmt.Print(start)    //display time
	itemCode := itmNo
	qty := q

	// print display on screen
	fmt.Printf("Item added %s qty=%d\n", itemCode, qty)
	start = time.Now() //start what time
	fmt.Print(start)   //display time

	my_res := "OK : checkin"

	result <- my_res
} //. End checkin

func checkout(itmNo string, q int16, result chan<- string) {
	start := time.Now() //start what time
	fmt.Print(start)    //display time
	itemCode := itmNo
	qty := q

	// print display on screen
	fmt.Printf("Item had checkout %s qty=%d\n", itemCode, qty)
	start = time.Now() //start what time
	fmt.Print(start)   //display time

	my_res := "OK : checkout"

	result <- my_res
} //. End checkin
