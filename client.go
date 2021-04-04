package main

import (
	"bufio"
	// "fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"

	// "time"
	// "os"
	"strings"
)

func main() {
	//Connection to server
	con, err := net.Dial("tcp", "0.0.0.0:9999")
	//Error condition
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()

	//Let server read what we sent
	serverReader := bufio.NewReader(con)

	//Send automatically 100 times
	for i := 0; i < 100; i++ {
		// Waiting for the client request
		clientRequest := strconv.Itoa(rand.Intn(4) + 1)
		//Write to server, about what client want
		if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
			log.Printf("failed to send the client request: %v\n", err)
		}

		// Waiting for the server response
		serverResponse, err := serverReader.ReadString('\n')

		switch err {
		//If no error, print response from server
		case nil:
			log.Println(strings.TrimSpace(serverResponse))
		//End of file
		case io.EOF:
			log.Println("server closed the connection")
			return
		//If don't have anyhting
		default:
			log.Printf("server error: %v\n", err)
			return
		}
	}
}
