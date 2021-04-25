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
	// "bytes"
	"os"
	"strings"
	"github.com/wcharczuk/go-chart" //exposes "chart"
)

func main() {

	x:= []float64{}
	y:= []float64{}
	start_whole := time.Now()
	var n float64
	n  = 1500
	start_n := time.Now()
	end := make(chan int, int(n))
	for i := 0; i < int(n); i++ {
		
		// Waiting for the client request
		go client(end,start_whole)

	}
	for i := 0; i < int(n); i++ {
		<-end
	}
	x = append(x, n)
	y = append(y, float64(time.Since(start_n)))

	fmt.Printf("Total time: %v\n", time.Since(start_whole))
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				// Style: chart.Style{
				// 	StrokeWidth:      chart.Disabled,
				// 	DotWidth:         5,
				// },
				XValues: x,
				YValues: y,
			},
		},
	}
	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func client(end chan int, start time.Time) {
	con, err := net.Dial("tcp", "127.0.0.3:8888")
		if err != nil {
			log.Fatalln(err)
		}
		defer con.Close()
		d:=50
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