package main

import (
	"fmt"

	"bookstore/bookstore"

	"github.com/kr/pretty"
)

func main() {
	db, err := bookstore.ConnectDB("dev", "playground")
	if err != nil {
		panic(err)
	}

	store := bookstore.NewBookstore(db)
	books := store.ReadBooks()

	fmt.Printf("dump: % #v\n", pretty.Formatter(books))
	fmt.Printf("dump: %s\n", books[0].ISBN)
}
