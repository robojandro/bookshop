package service

import (
	"fmt"

	"bookshop/books"

	uuid "github.com/satori/go.uuid"
)

// BookDataStore provides an interface for interacting with the BookDataStore.
type BookDataStore interface {
	DeleteBooks(ids ...string) error
	ReadBookByISBN(isbn string) (books.Book, error)
	ReadBooks() ([]books.Book, error)
	UpsertBooks(books []books.Book) error
}

// SVC is an interface that fulfills bookshop service calls.
type SVC interface {
	AddBook(title, isbn string) (books.Book, error)
	ListBooks() ([]books.Book, error)
}

// Service is a wrapper for the bookshop service business logic.
type Service struct {
	bookStore BookDataStore
}

// NewService returns a Service type value.
func NewService(bs BookDataStore) Service {
	return Service{
		bookStore: bs,
	}
}

// AddBook add a book from the given title and isbn if the isbn does not already exist.
func (s *Service) AddBook(title, isbn string) (books.Book, error) {
	//make sure book doesn't already exist
	extant, err := s.bookStore.ReadBookByISBN(isbn)
	if err != nil {
		return books.Book{}, err
	}
	if isbn == string(extant.ISBN) {
		return books.Book{}, NewErrDuplicate(title)
	}

	bks := []books.Book{
		{
			ID:    fmt.Sprintf("%s", uuid.NewV4()),
			Title: title,
			ISBN:  books.ISBN(isbn),
		},
	}
	if err := s.bookStore.UpsertBooks(bks); err != nil {
		return books.Book{}, err
	}

	return bks[0], nil
}

func (s *Service) ListBooks() ([]books.Book, error) {
	return s.bookStore.ReadBooks()
}

func (s *Service) RemoveBooks(ids ...string) error {
	return s.bookStore.DeleteBooks(ids...)
}

func (s *Service) UpdateBook(bk books.Book) error {
	return s.bookStore.UpsertBooks([]books.Book{bk})
}
