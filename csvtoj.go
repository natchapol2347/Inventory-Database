package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"math/rand"
)
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}
func main(){
empData := [][]string{
	{"Id", "Name", "Quantity","ExpDate"},
}


for i:=1;i<=100;i++{
	s := []string{strconv.Itoa(i),randSeq(7),"30","22032021"}
	empData = append(empData,s)
}



csvFile, err := os.Create("input.csv")

if err != nil {
	log.Fatalf("failed creating file: %s", err)
}

csvwriter := csv.NewWriter(csvFile)
 
for _, empRow := range empData {
	_ = csvwriter.Write(empRow)
}
csvwriter.Flush()
csvFile.Close()
}