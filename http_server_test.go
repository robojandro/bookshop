package main

import (
	"bookshop/books"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPServer(t *testing.T) {
	t.Run("books", func(t *testing.T) {
		t.Run("GET", func(t *testing.T) {
			t.Run("happy", func(t *testing.T) {
				mockBooks = []books.Book{
					{
						ID:    "abc01",
						Title: "titleA",
						ISBN:  "9783161484100",
					},
					{
						ID:    "def02",
						Title: "titleB",
						ISBN:  "2222222222222",
					},
				}

				req, err := http.NewRequest("GET", "/books", nil)
				if err != nil {
					t.Fatal(err)
				}
				resp := makeRequest(req)

				require.Equal(t, resp.Code, http.StatusOK)

				var content []books.Book
				err = json.NewDecoder(resp.Body).Decode(&content)
				require.NoError(t, err)
				assert.Equal(t, mockBooks, content)
			})

			t.Run("service error", func(t *testing.T) {
				mockBooks = nil
				mockBooksErr = errors.New("service error")

				req, err := http.NewRequest("GET", "/books", nil)
				if err != nil {
					t.Fatal(err)
				}
				resp := makeRequest(req)

				assert.Equal(t, resp.Body.String(), "service error\n")
				assert.Equal(t, resp.Code, http.StatusInternalServerError)
			})
		})
	})
}

func makeRequest(req *http.Request) *httptest.ResponseRecorder {
	svc := &mockService{}
	httpServer := NewHTTPServer(svc)
	resp := httptest.NewRecorder()
	httpServer.ServeHTTP(resp, req)
	return resp
}

type mockService struct{}

var mockBooks []books.Book
var mockBooksErr error

func (m *mockService) ListBooks() ([]books.Book, error) {
	return mockBooks, mockBooksErr
}
