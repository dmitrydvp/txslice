package main

import (
	"fmt"

	"github.com/d1m4ek1/txslice"
)

type Values struct {
	Val int
}

func main() {
	tx := txslice.New([]*Values{}, txslice.Config{})

	// Push elements
	tx.Push(&Values{1})
	tx.Push(&Values{2})

	fmt.Println("First element:", *tx.FirstElement()) // {1}

	// Commit example
	tx.Commit()

	// Transaction with rollback
	tx.Push(&Values{3})
	tx.Rollback() // rollback last transaction

	fmt.Println("Length after rollback:", tx.Len()) // 2
}
