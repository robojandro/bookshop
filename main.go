package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"bookshop/books"
	"bookshop/datastore"
	"bookshop/service"
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

	bookStore := books.NewBookStore(data)
	service := service.NewService(&bookStore)
	httpServer := NewHTTPServer(&service)

	srv := &http.Server{
		Handler:      httpServer,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
