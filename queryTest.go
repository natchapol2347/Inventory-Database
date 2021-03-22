package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"net"
	"bufio"
	"io"
	"log"
	"strings"
)

func main() {

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err := sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
	listener, err := net.Listen("tcp","0.0.0.0:8000")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for{
		
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleClientRequest(con)
		// if there is an error opening the connection, handle it
		if err != nil {
			panic(err.Error())
		}

		// defer the close till after the main function has finished
		// executing
		defer db.Close()
		fmt.Println("Successfully connect to database")
		// perform a db.Query insert
		insert, err := db.Query("INSERT INTO test VALUES(1)")

		// if there is an error inserting, handle it
		if err != nil {
			panic(err.Error())
		}
		// be careful deferring Queries if you are using transactions
		defer insert.Close()
		fmt.Println("Inserted complete")
	}
	

}

func handleClientRequest(con net.Conn) {
	defer con.Close()

	clientReader := bufio.NewReader(con)

	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if clientRequest == ":QUIT" {
				log.Println("client requested server to close the connection so closing")
				return
			} else {
				log.Println(clientRequest)
			}
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}

		// Responding to the client request
		if _, err = con.Write([]byte("GOT IT!\n")); err != nil {
			log.Printf("failed to respond to client: %v\n", err)
		}
	}
}
