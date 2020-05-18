package service_test

import (
	"errors"
	"testing"
	"time"

	"bookshop/authors"
	"bookshop/books"
	"bookshop/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	dt, err := time.Parse(authors.DateParsingFormat, "1970-01-01")
	require.NoError(t, err)

	t.Run("GetAuthor", func(t *testing.T) {
		t.Run("happy", func(t *testing.T) {
			mockAuthStore := &mockAuthorStore{}
			srv := service.NewService(mockAuthStore, nil)

			mockAuthErr = nil
			mockAuth = authors.Author{
				ID:         "auth01",
				FirstName:  "First",
				MiddleName: "Middle",
				LastName:   "Last",
				DOB:        &dt,
				Books: []books.Book{
					{
						ID:    "abc01",
						Title: "titleA",
						ISBN:  "9783161484100",
					},
				},
			}
			_, err := srv.GetAuthor("auth01")
			assert.NoError(t, err)
		})
	})

	t.Run("AddBook", func(t *testing.T) {
		t.Run("happy", func(t *testing.T) {
			mockBkStore := &mockBookStore{}
			srv := service.NewService(nil, mockBkStore)
			mockBooks = nil

			_, err := srv.AddBook("titleA", "9783161484100")
			assert.NoError(t, err)
		})

		t.Run("already exists", func(t *testing.T) {
			mockBkStore := &mockBookStore{}
			srv := service.NewService(nil, mockBkStore)
			mockBooksErr = nil
			mockBook = books.Book{
				ID:    "abc01",
				Title: "titleA",
				ISBN:  "9783161484100",
			}

			results, err := srv.AddBook("titleA", "9783161484100")
			assert.Error(t, err)
			assert.Equal(t, results, books.Book{})
		})

		t.Run("datastore error", func(t *testing.T) {
			mockBkStore := &mockBookStore{}
			srv := service.NewService(nil, mockBkStore)
			mockBooks = nil
			mockBooksErr = errors.New("datastore error")

			results, err := srv.AddBook("titleA", "9783161484100")
			assert.Error(t, err)
			assert.Equal(t, results, books.Book{})
		})
	})

	t.Run("ListBooks", func(t *testing.T) {
		t.Run("happy", func(t *testing.T) {
			mockBkStore := &mockBookStore{}
			srv := service.NewService(nil, mockBkStore)
			mockBooksErr = nil
			mockBooks = []books.Book{
				{
					ID:    "abc01",
					Title: "titleA",
					ISBN:  "9783161484100",
				},
			}

			results, err := srv.ListBooks()
			assert.NoError(t, err)
			assert.Len(t, results, 1)
			assert.Equal(t, mockBooks, results)
		})

		t.Run("datastore error", func(t *testing.T) {
			mockBkStore := &mockBookStore{}
			srv := service.NewService(nil, mockBkStore)
			mockBooks = nil
			mockBooksErr = errors.New("datastore error")

			results, err := srv.ListBooks()
			assert.Error(t, err)
			assert.Nil(t, results)
		})
	})

	t.Run("RemoveBooks", func(t *testing.T) {
		t.Run("happy", func(t *testing.T) {
			mockBkStore := &mockBookStore{}
			srv := service.NewService(nil, mockBkStore)
			mockBooksErr = nil
			err := srv.RemoveBooks("abc01", "def02")
			assert.NoError(t, err)
		})

		t.Run("datastore error", func(t *testing.T) {
			mockBkStore := &mockBookStore{}
			srv := service.NewService(nil, mockBkStore)
			mockBooksErr = errors.New("datastore error")
			err := srv.RemoveBooks("abc01", "def02")
			assert.Error(t, err)
		})
	})

	t.Run("UpdateBook", func(t *testing.T) {
		mockBook = books.Book{
			ID:    "abc01",
			Title: "titleA",
			ISBN:  "9783161484100",
		}
		t.Run("happy", func(t *testing.T) {
			mockBkStore := &mockBookStore{}
			srv := service.NewService(nil, mockBkStore)
			mockBooksErr = nil
			err := srv.UpdateBook(mockBook)
			assert.NoError(t, err)
		})

		t.Run("datastore error", func(t *testing.T) {
			mockBkStore := &mockBookStore{}
			srv := service.NewService(nil, mockBkStore)
			mockBooksErr = errors.New("datastore error")
			err := srv.UpdateBook(mockBook)
			assert.Error(t, err)
		})
	})
}
