package service_test

import (
	"bookshop/authors"
	"bookshop/books"
)

var mockAuth authors.Author
var mockAuths []authors.Author
var mockAuthErr error

type mockAuthorStore struct{}

func (m *mockAuthorStore) ReadAuthors() ([]authors.Author, error) {
	return mockAuths, mockAuthErr
}

func (m *mockAuthorStore) ReadAuthorAndBooks(id string) (authors.Author, error) {
	return mockAuth, mockAuthErr
}

var mockBook books.Book
var mockBooksErr error
var mockBooks []books.Book

type mockBookStore struct{}

func (m *mockBookStore) DeleteBooks(ids ...string) error {
	return mockBooksErr
}

func (m *mockBookStore) ReadBookByISBN(isbn string) (books.Book, error) {
	return mockBook, mockBooksErr
}

func (m *mockBookStore) ReadBooks() ([]books.Book, error) {
	return mockBooks, mockBooksErr
}

func (m *mockBookStore) UpsertBooks(books []books.Book) error {
	return mockBooksErr
}
