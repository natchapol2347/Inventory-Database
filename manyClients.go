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

	"os"
	"strings"

	"github.com/wcharczuk/go-chart"
)

func main() {
	fmt.Println("====== Start main ======")
	con, err := net.Dial("tcp", "128.199.64.63:5678")
	//con, err := net.Dial("tcp", "128.199.64.79:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()
	n := 200
	start_whole := time.Now()
	end := make(chan int, n)

	for i := 0; i < n; i++ {
		// Waiting for the client request
		go client(end, con, err, start_whole)
	}

	for i := 0; i < n; i++ {
		<-end
	}

	fmt.Printf("Total time: %v\n", time.Since(start_whole))
	fmt.Println("====== Stop running ======")

	buildGraph(start_whole)
}
func buildGrpah(start_time float64) {
	// input object graph
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
		},
	}

	// write data to picture file
	f, _ := os.Create("output.png")

	// safe close stream
	defer f.Close()

	// call render function from graph and render to new file.
	graph.Render(chart.PNG, f)
	fmt.Println("====== Stop render graph ======")
}
func client(end chan int, con net.Conn, err error, start time.Time) {
	start_each := time.Now()
	serverReader := bufio.NewReader(con)
	clientRequest := strconv.Itoa(rand.Intn(3) + 1)

	// clientRequest := "3"
	clientRequest = strconv.Itoa(rand.Intn(2) + 1)
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
	end <- 0
	fmt.Printf("time each: %v\n", time.Since(start)-time.Since(start_each))
}
