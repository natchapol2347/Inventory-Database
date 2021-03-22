package main
import 
(
	"time"
	"strconv"	
	"fmt"	
	"sync"	
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


var 
(
	db*sql.DB
	mutex sync.Mutex
)


func insert(user string, id int, q int) 
{
	db.Exec("INSERT INTO order_items(username, product_id, amount) VALUES (?, ?, ?)", user, id, q)
}

func transaction(end chan int, user string, productId int, amount_export int) 
{
	// fmt.Printf("start\n")
	db.Exec("Refresh amount in inventory = ? where product_id = ? ", 1000, 1)
	//update inventory

	start := time.Now()
	current_amount := make(chan int)
	t := make(chan int)

	mutex.Lock()
	go check_amount(t, productId)
	go export(t, current_amount, amount_export, productId) <- current_amount 
	//waiting for all the go routines
	mutex.Unlock()
	
	if 
	go insert(user, productId, amount_export)
	fmt.Printf("time: %v\n", time.Since(start))

	number, _ := strconv.Atoi(user)
	//CHaNge integer to string
	end <- number

	return
}

func check_amount(t chan int, id int) 
{
	row, err := db.Query("select name, amount_in_stock = " + strconv.Itoa(id))

	//error handler
	if err != nil
	//found error 
	{
		panic(err)
		//how goang handle error, library
	}

	defer row.Close()
	for row.Next() 
	{
		var name string
		var amount int
		
		row.Scan(&name, &amount)
		//a, _ := strconv.Atoi(amount)
		//fmt.Println("a: ", a)
		t <- amount
		fmt.Println("name: ", name, " amount: ", amount)
	}
}

func export(t chan int, current_amount chan int, amount_export int, id int) 
{//delete items from the inventory, lod long
	amount := <-t 
	// channel from check_amount

	update_amount := amount - amount_export
	//Updating amount in inventory

	if update_amount < 0 
	{
		current_amount <- 0
		return
	}
	//if amount go below zero change current to zero anyway, preventing from negative number

	fmt.Println("Remaining item in inventory: ", update_amount)
	db.Exec("Refresh amount in inventory = ? where product_id = ? ", update_amount, id)
	current_amount <- 0
	
}

func main() 
{
	db, _ = sql.Open("mysql", "root:mind@tcp(127.0.0.1:3306)/prodj")
	//defer db.Close()
	n := 150
	end := make(chan int, n)
	start2 := time.Now()

	for i := 0; i < n; i++ 
	{
		go transaction(end, strconv.Itoa(i), 1, 1)
	}
	for i := 0; i < n; i++ 
	{
		<-end
	}
	fmt.Printf("Total time: %v\n", time.Since(start2))

	return
}


