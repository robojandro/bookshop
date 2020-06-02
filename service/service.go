package service

import (
	"bookshop/authors"
	"bookshop/books"

	uuid "github.com/satori/go.uuid"
)

// AuthorDataStore provides an interface for interacting with the AuthorDataStore.
type AuthorDataStore interface {
	DeleteAuthor(id string) error
	ReadAuthors() ([]authors.Author, error)
	ReadAuthorAndBooks(id string) (authors.Author, error)
}

// BookDataStore provides an interface for interacting with the BookDataStore.
type BookDataStore interface {
	DeleteBooks(ids ...string) error
	ReadBookByISBN(isbn string) (books.Book, error)
	ReadBooks() ([]books.Book, error)
	UpsertBooks(books []books.Book) error
}

// SVC is an interface that fulfills bookshop service calls.
type SVC interface {
	GetAuthor(id string) (authors.Author, error)
	ListAuthors() ([]authors.Author, error)
	RemoveAuthor(id string) error

	AddBook(title, isbn string) (books.Book, error)
	RemoveBooks(ids ...string) error
	ListBooks() ([]books.Book, error)
	UpdateBook(bk books.Book) error
}

// Service is a wrapper for the bookshop service business logic.
type Service struct {
	authStore AuthorDataStore
	bookStore BookDataStore
}

// NewService returns a Service type value.
func NewService(as AuthorDataStore, bs BookDataStore) Service {
	return Service{
		authStore: as,
		bookStore: bs,
	}
}

// GetAuthor will return the details for an author by the given id including
// their list of books.
func (s *Service) GetAuthor(id string) (authors.Author, error) {
	return s.authStore.ReadAuthorAndBooks(id)
}

// ListAuthors will return a list of all authors sorted by last name (ascending).
func (s *Service) ListAuthors() ([]authors.Author, error) {
	return s.authStore.ReadAuthors()
}

// RemoveAuthor will delete the author from the datastore.
func (s *Service) RemoveAuthor(id string) error {
	return s.authStore.DeleteAuthor(id)
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
			ID:    uuid.NewV4().String(),
			Title: title,
			ISBN:  books.ISBN(isbn),
		},
	}
	if err := s.bookStore.UpsertBooks(bks); err != nil {
		return books.Book{}, err
	}

	return bks[0], nil
}

// ListBooks will return a list of books.Book.
func (s *Service) ListBooks() ([]books.Book, error) {
	return s.bookStore.ReadBooks()
}

// RemoveBooks will remove a list of books by the given book.Book IDs.
func (s *Service) RemoveBooks(ids ...string) error {
	return s.bookStore.DeleteBooks(ids...)
}

// UpdateBook will update the given book with the information from the request.
func (s *Service) UpdateBook(bk books.Book) error {
	return s.bookStore.UpsertBooks([]books.Book{bk})
}
