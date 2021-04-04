package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

// Person is an object to collect person data
type ProductItem struct {
	ItemCode string
	ItemName string
	qty  int
}



func main() {
	arguments := os.Args
        if len(arguments) == 1 {
                fmt.Println("Please provide host:port.")
                return
        }

        CONNECT := arguments[1]
        c, err := net.Dial("tcp", CONNECT)
        if err != nil {
                fmt.Println(err)
                return
        }


		// call menu 
		go selectMenu()



        for {
				// display menu and wait to display on screen
				go selecrtMenu()

                reader := bufio.NewReader(os.Stdin)
                fmt.Print(">> ")

                text, _ := reader.ReadString('\n')
                
				// show out put from server
				fmt.Fprintf(c, text+"\n")

                message, _ := bufio.NewReader(c).ReadString('\n')
                fmt.Print("->: " + message)


				
                if strings.TrimSpace(string(text)) == "STOP" {
                        fmt.Println("TCP client exiting...")
                        return
                }
        }

} // .end main

func int selectMenu() {
	//Create menu for input
	fmt.Print("Menu Program \n")
	fmt.Print("1. Input Item\n")
	fmt.Print("2. Checkout Item\n")
	fmt.Print("3. Show Item\n")
	
	fmt.Print("-------------\n")
	fmt.Print("Pls, select menu from above = ")
	d:= input.Scan()
	return d
}

func checkin(itmNo string, qty int16) {
	//manual enter
	reader := bufio.NewReader(os.Stdin)

	//search for items in inventory
	fmt.Print("Enter Item Code: ")
    itemCode, _ := reader.ReadString('\n')
	
    
	fmt.Print("Enter Qty: ")
    qty, _ := reader.ReadString('\n')

	go connectToCheckin()

	// print display on screen
	fmt.Printf("Item added %s\n", itemCode)

} //. End checkin

func connectToChekcin(){

}

func connectToChekcout(){
	
}

func checkout(itmNo string, qty int16) {
	//
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Item Code: ")
    itemCode, _ := reader.ReadString('\n')
	
	fmt.Print("Enter Item Name: ")
    itemName, _ := reader.ReadString('\n')
    
	fmt.Print("Enter Qty: ")
    qty, _ := reader.ReadString('\n')


	go connectToCheckout()

	// print display on screen
	fmt.Printf("Item got %s\n", itemCode)

} //. End checkout

func currentStock(jsonString string) {
	
} // .end

// do sum process when connected
func handleConnection(c net.Conn) {
	// print display on screen
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			// 
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}


		while () {

		}


		result := "test\n"
		c.Write([]byte(string(result)))
	}
	c.Close()
} // .End handleConnection
