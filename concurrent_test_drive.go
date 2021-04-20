package main

import (
	"fmt"
)

func main() {
	fmt.Println("====== Start main ======")

	// declare channel
	my_res := make(chan string)
	// test concurrent
	for i := 0; i < 10; i++ {
		// call function
		go checkin("123", 2, my_res)
	} //. End for

	// display my_res on screen
	fmt.Println(<-my_res)

	// display text on screen
	fmt.Println("====== Start main ======")
} // .End function
