package main

import (
	"bufio"
	// "fmt"
	"math/rand"
	"strconv"
	"io"
	"log"
	"net"
	// "time"
	// "os"
	"strings"
)


//Old file of client
func main() {
	con, err := net.Dial("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()

	serverReader := bufio.NewReader(con)

	for i:=0; i<100; i++{
		// Waiting for the client request
		clientRequest:= strconv.Itoa(rand.Intn(4)+1) //Randomize 5 numbers
			if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
				log.Printf("failed to send the client request: %v\n", err)
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
}
