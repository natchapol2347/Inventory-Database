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

        for {
				// display menu and wait to display on screen
				go selecrtMenu()

				start2 := time.Now()

				for i := 0; i < n; i++ {
					go checkin(endin, strconv.Itoa(i), 1, 10)
				}

				for i := 0; i < n; i++ {
					go checkout(endin, strconv.Itoa(i), 1, 10)
				}

                reader := bufio.NewReader(os.Stdin)
                //fmt.Print(">> ")

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
	//
	reader := bufio.NewReader(os.Stdin)

	//fmt.Print("Enter Item Code: ")
    itemCode, _ := reader.ReadString('\n')
	
    
	//fmt.Print("Enter Qty: ")
    qty, _ := reader.ReadString('\n')

	go connectToCheckin()

	// print display on screen
	fmt.Printf("Item added %s\n", itemCode)

} //. End checkin

func connectToChekcinn chan string, e chan int, quantity int, id int, name string, con net.Conn) {
	defer con.Close()
 
	clientReader := bufio.NewReader(con)

	product := <-n
	expdate := <-e

	// Responding to the client request
	_, err = con.Write([]byte("1\n"))
	
}

func connectToChekcout(n chan string, e chan int, quantity int, id int, name string, con net.Conn) {
	defer con.Close()
 
	clientReader := bufio.NewReader(con)
	product := <-n
	expdate := <-e
	
	// Responding to the client request
	_, err = con.Write([]byte("2\n"))
}
	

func checkout(itmNo string, qty int16) {
	//
	reader := bufio.NewReader(os.Stdin)

	//fmt.Print("Enter Item Code: ")
    itemCode, _ := reader.ReadString('\n')
	
	//fmt.Print("Enter Item Name: ")
    itemName, _ := reader.ReadString('\n')
    
	//fmt.Print("Enter Qty: ")
    qty, _ := reader.ReadString('\n')


	go connectToCheckout()

	// print display on screen
	// fmt.Printf("Item got %s\n", itemCode)

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

		result := "test\n"
		c.Write([]byte(string(result)))
	}
	c.Close()
} // .End handleConnection
