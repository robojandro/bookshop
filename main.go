package main

import (
	"fmt"
	"log"
	"os"

	"bookshop/books"
	"bookshop/datastore"

	"github.com/kr/pretty"
)

func main() {
	var dbUser, dbPass, dbName string
	if dbUser = os.Getenv("bookshop_dbuser"); dbUser == "" {
		log.Fatal("missing env variable bookshop_dbuser")
	}
	if dbPass = os.Getenv("bookshop_dbpass"); dbPass == "" {
		log.Fatal("missing env variable bookshop_dbpass")
	}
	if dbName = os.Getenv("bookshop_dbname"); dbName == "" {
		log.Fatal("missing env variable bookshop_dbname")
	}

	data, err := datastore.ConnectDB(dbUser, dbPass, dbName)
	if err != nil {
		panic(err)
	}

	store := books.NewBookStore(data)
	books, err := store.ReadBooks()
	if err != nil {
		panic(err)
	}

	fmt.Printf("books: % #v\n", pretty.Formatter(books))
	//fmt.Printf("ISBN: %s\n", books[0].ISBN)
}
