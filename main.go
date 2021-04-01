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

func main2() {
	con, err := net.Dial("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()

	clientReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(con)

	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
				log.Printf("failed to send the client request: %v\n", err)
			}
		case io.EOF:
			log.Println("client closed the connection")
			return
		default:
			log.Printf("client error: %v\n", err)
			return
		}

		// Waiting for the server response
		serverResponse, err := serverReader.ReadString('\n')

		switch err {
		case nil:
			log.Println(strings.TrimSpace(serverResponse))
		case io.EOF:
			log.Println("server closed the connection")
			return
		default:
			log.Printf("server error: %v\n", err)
			return
		}
	}

} // .End Main

func main() {
	arguments := os.Args // detect parameter

	// count argument size of array
	if len(arguments) == 1 {
		// display text to screen
		fmt.Println("Please provide a port number!")
		// stop program
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		// display text to screen
		fmt.Println(err)
		return
	}

	defer l.Close()

	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}

} // .end main

func checkin(itmNo string, qty int16) {
	//
	reader := bufio.NewReader(os.Stdin)

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
