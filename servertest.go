package main
 
import (
	"bufio"
	"database/sql"
	"io"
	"log"
	"net"
	"strings"
	"funcs"
	_ "github.com/go-sql-driver/mysql"
)
 
var (
	db    *sql.DB
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
		number:=0
		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if clientRequest == "QUIT" {
				log.Println("client requested server to close the connection so closing")
				return
			}else if clientRequest == "1"{
				log.Println("Insert items")
				message="Insert items"
				number=1
			}else if clientRequest == "2"{
				log.Println("Remove items")
				message="Remove items"
				number=2
			}else if clientRequest == "3"{
				log.Println("Check current stock")
				message="Check current stock"
				number=3
			}else if clientRequest == "4"{
				log.Println("Check record for insert")
				message="Check record for insert"
				number=4
			}else if clientRequest == "5"{
				log.Println("Check record for remove")
				message="Check record for remove"
				number=5
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
		
		if number == 1{
			db, _ = sql.Open("mysql", "ohm:!Bruno555@tcp(127.0.0.1:3306)/inventory")
			endin := make(chan int)
			funcs.Going_in(endin,"user", 1, 1)
			<-endin
		}else if number == 2{
			//run going_out
		}else if number == 3{
			//run show_current
		}else if number == 4{
			//run show_record_in
		}else if number == 5{
			//run show_record_out
		}

	}
}
