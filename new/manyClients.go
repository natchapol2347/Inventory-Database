package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	// "os"
	"strings"
)

func main() {
	

	n := 50
	start_whole := time.Now()
	end := make(chan int, n)
	for i := 0; i < n; i++ {
		
		// Waiting for the client request
		go client(end,start_whole)

	}
	for i := 0; i < n; i++ {
		<-end
	}
	fmt.Printf("Total time: %v\n", time.Since(start_whole))
}

func client(end chan int, start time.Time) {
	con, err := net.Dial("tcp", "0.0.0.0:5678")
		if err != nil {
			log.Fatalln(err)
		}
		defer con.Close()
		d:=100
		for i:=0;i<d;i++{
			start_each := time.Now()
			serverReader := bufio.NewReader(con)
			clientRequest := strconv.Itoa(rand.Intn(2) + 1)
			// clientRequest := "4"
			if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
				log.Printf("failed to send the client request: %v\n", err)
			}

			// Waiting for the server response
			serverResponse, err := serverReader.ReadString('.')

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
			fmt.Printf("time each: %v\n", time.Since(start)-time.Since(start_each))
		}
		end <- 0
		
		
}