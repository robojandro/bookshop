package main

import (
	"fmt"
	"log"
	"os"

	"bookstore/bookstore"

	"github.com/kr/pretty"
)

func main() {
	var dbUser, dbPass, dbName string
	if dbUser = os.Getenv("bookstore_dbuser"); dbUser == "" {
		log.Fatal("missing env variable bookstore_dbuser")
	}
	if dbPass = os.Getenv("bookstore_dbpass"); dbPass == "" {
		log.Fatal("missing env variable bookstore_dbpass")
	}
	if dbName = os.Getenv("bookstore_dbname"); dbName == "" {
		log.Fatal("missing env variable bookstore_dbname")
	}

	db, err := bookstore.ConnectDB(dbUser, dbPass, dbName)
	if err != nil {
		panic(err)
	}

	store := bookstore.NewBookstore(db)
	books, err := store.ReadBooks()
	if err != nil {
		panic(err)
	}

	fmt.Printf("books: % #v\n", pretty.Formatter(books))
	//fmt.Printf("ISBN: %s\n", books[0].ISBN)
}
