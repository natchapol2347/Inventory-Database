package main
 
import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)
 
func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
 
		// If you want, you can increment a counter here and inject to handleClientRequest below as client identifier
		go handleClientRequest(con)
	}
}
 
func handleClientRequest(con net.Conn) {
	defer con.Close()
 
	clientReader := bufio.NewReader(con)
 
	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')
		message:="Please provide numbers 1-5"
		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if clientRequest == "QUIT" {
				log.Println("client requested server to close the connection so closing")
				return
			}else if clientRequest == "1"{
				log.Println("Insert items")
				message="Insert items"
			}else if clientRequest == "2"{
				log.Println("Remove items")
				message="Remove items"
			}else if clientRequest == "3"{
				log.Println("Check current stock")
				message="Check current stock"
			}else if clientRequest == "4"{
				log.Println("Check record for insert")
				message="Check record for insert"
			}else if clientRequest == "5"{
				log.Println("Check record for remove")
				message="Check record for remove"
			}else {
				log.Println("Please provide numbers 1-5")
			}
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}
 
		// Responding to the client request
		_, err = con.Write([]byte(message+"\n"))
		if err != nil {
			log.Printf("failed to respond to client: %v\n", err)
		}
	}
}
