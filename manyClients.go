package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
	// "os"
	"strings"
)
 
func main() {
	con, err := net.Dial("tcp", "0.0.0.0:8000")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()
 
	serverReader := bufio.NewReader(con)
	count:=0
	start := time.Now()
	for i:=0;i<20000;i++{
		// Waiting for the client request
		clientRequest := "banana"
 
		switch err {
		case nil:
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
			count = count+1
			fmt.Println(count)
		case io.EOF:
			log.Println("server closed the connection")
			return
		default:
			log.Printf("server error: %v\n", err)
			return
		}
	}
	fmt.Printf("Total time: %v\n", time.Since(start))
}