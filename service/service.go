package service

import (
	"fmt"

	"bookshop/books"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// Store provides an interface for interacting with groups.
type BookDataStore interface {
	DeleteBooks(ids ...string) error
	InsertBooks(books []books.Book) error
	ReadBookByISBN(isbn string) (books.Book, error)
	ReadBooks() ([]books.Book, error)
	UpsertBooks(books []books.Book) error
}

// Service is a wrapper for the custom field fielddata business logic.
type Service struct {
	bookStore BookDataStore
}

func NewService(bs BookDataStore) Service {
	return Service{
		bookStore: bs,
	}
}

// AddBooks
func (s *Service) AddBook(title, isbn string) (books.Book, error) {
	//make sure book doesn't already exist
	extant, err := s.bookStore.ReadBookByISBN(isbn)
	if err != nil {
		return books.Book{}, err
	}
	if isbn == string(extant.ISBN) {
		return books.Book{}, errors.Errorf("isbn already exists with title: %s", extant.Title)
	}

	bks := []books.Book{
		{
			ID:    fmt.Sprintf("%s", uuid.NewV4()),
			Title: title,
			ISBN:  books.ISBN(isbn),
		},
	}
	if err := s.bookStore.InsertBooks(bks); err != nil {
		return books.Book{}, err
	}

	return books.Book{}, nil
}

func (s *Service) AllBooks() ([]books.Book, error) {
	return s.bookStore.ReadBooks()
}
