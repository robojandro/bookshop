package service_test

import (
	"testing"

	"bookshop/books"
	"bookshop/service"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	t.Run("AddBook", func(t *testing.T) {
		t.Run("happy", func(t *testing.T) {
			mockBkStore := &mockBookStore{}
			srv := service.NewService(mockBkStore)
			mockBooks = nil

			_, err := srv.AddBook("titleA", "9783161484100")
			assert.NoError(t, err)
		})

		t.Run("already exists", func(t *testing.T) {
			mockBkStore := &mockBookStore{}
			srv := service.NewService(mockBkStore)
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
			srv := service.NewService(mockBkStore)
			mockBooks = nil
			mockBooksErr = errors.New("datastore error")

			results, err := srv.AddBook("titleA", "9783161484100")
			assert.Error(t, err)
			assert.Equal(t, results, books.Book{})
		})
	})

	t.Run("AllBooks", func(t *testing.T) {
		t.Run("happy", func(t *testing.T) {
			mockBkStore := &mockBookStore{}
			srv := service.NewService(mockBkStore)
			mockBooksErr = nil
			mockBooks = []books.Book{
				{
					ID:    "abc01",
					Title: "titleA",
					ISBN:  "9783161484100",
				},
			}

			results, err := srv.AllBooks()
			assert.NoError(t, err)
			assert.Len(t, results, 1)
			assert.Equal(t, mockBooks, results)
		})

		t.Run("datastore error", func(t *testing.T) {
			mockBkStore := &mockBookStore{}
			srv := service.NewService(mockBkStore)
			mockBooks = nil
			mockBooksErr = errors.New("datastore error")

			results, err := srv.AllBooks()
			assert.Error(t, err)
			assert.Nil(t, results)
		})
	})
}
