package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bookshop/books"
	"bookshop/service"

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

				resp := makeRequest(t, "GET", "/books", "")

				require.Equal(t, http.StatusOK, resp.Code)

				var content []books.Book
				err := json.NewDecoder(resp.Body).Decode(&content)
				require.NoError(t, err)
				assert.Equal(t, mockBooks, content)
			})

			t.Run("service error", func(t *testing.T) {
				mockBooks = nil
				mockBooksErr = errors.New("service error")

				resp := makeRequest(t, "GET", "/books", "")

				assert.Equal(t, "service error\n", resp.Body.String())
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			})
		})

		t.Run("DELETE", func(t *testing.T) {
			t.Run("happy", func(t *testing.T) {
				mockBooksErr = nil
				resp := makeRequest(t, "DELETE", "/books/abc01", "")
				require.Equal(t, http.StatusAccepted, resp.Code)
			})

			t.Run("empty id", func(t *testing.T) {
				mockBooksErr = errors.New("book_id cannot be blank")
				resp := makeRequest(t, "DELETE", "/books/ ", "")
				require.Equal(t, http.StatusBadRequest, resp.Code)
			})

			t.Run("service error", func(t *testing.T) {
				mockBooksErr = errors.New("service error")
				resp := makeRequest(t, "DELETE", "/books/abc01", "")
				require.Equal(t, http.StatusInternalServerError, resp.Code)
			})
		})

		t.Run("PATCH", func(t *testing.T) {
			t.Run("happy", func(t *testing.T) {
				mockBooksErr = nil
				resp := makeRequest(t, "PATCH", "/books", `{
				  "id": "abc01",
				  "title": "UPDATED TITLE",
				  "isbn": "999999999"
				}`)
				require.Equal(t, http.StatusAccepted, resp.Code)
			})

			t.Run("missing id", func(t *testing.T) {
				mockBooksErr = nil
				resp := makeRequest(t, "PATCH", "/books", `{
				  "title": "UPDATED TITLE",
				  "isbn": "999999999"
				}`)
				require.Equal(t, http.StatusBadRequest, resp.Code)
			})

			t.Run("service error", func(t *testing.T) {
				mockBooksErr = errors.New("service error")
				resp := makeRequest(t, "PATCH", "/books", `{
				  "id": "abc01",
				  "title": "UPDATED TITLE",
				  "isbn": "999999999"
				}`)
				require.Equal(t, http.StatusInternalServerError, resp.Code)
			})
		})

		t.Run("POST", func(t *testing.T) {
			t.Run("happy", func(t *testing.T) {
				mockBooksErr = nil
				mockBook = books.Book{
					ID:    "abc01",
					Title: "title3",
					ISBN:  "999999999",
				}

				resp := makeRequest(t, "POST", "/books", `{"title": "title03", "isbn": "999999999"}`)

				require.Equal(t, http.StatusCreated, resp.Code)

				var content books.Book
				err := json.NewDecoder(resp.Body).Decode(&content)
				require.NoError(t, err)
				assert.Equal(t, mockBook, content)
			})

			t.Run("duplicate", func(t *testing.T) {
				mockBook = books.Book{}
				mockBooksErr = service.NewErrDuplicate("title03")

				resp := makeRequest(t, "POST", "/books", `{"title": "title03", "isbn": "999999999"}`)
				require.Equal(t, http.StatusUnprocessableEntity, resp.Code)

				assert.Equal(t, mockBooksErr.Error(), strings.TrimSpace(resp.Body.String()))
			})

			t.Run("service error", func(t *testing.T) {
				mockBook = books.Book{}
				mockBooksErr = errors.New("service error")

				resp := makeRequest(t, "POST", "/books", `{"title": "title03", "isbn": "999999999"}`)
				require.Equal(t, http.StatusInternalServerError, resp.Code)
			})
		})
	})
}

func makeRequest(t *testing.T, kind, path, bodyStr string) *httptest.ResponseRecorder {
	bd := strings.NewReader(bodyStr)
	req, err := http.NewRequest(kind, path, bd)
	if err != nil {
		t.Fatal(err)
	}

	svc := &mockService{}
	httpServer := NewHTTPServer(svc)
	resp := httptest.NewRecorder()
	httpServer.ServeHTTP(resp, req)
	return resp
}

type mockService struct{}

var mockBook books.Book
var mockBooks []books.Book
var mockBooksErr error

func (m *mockService) AddBook(title, isbn string) (books.Book, error) {
	return mockBook, mockBooksErr
}

func (m *mockService) ListBooks() ([]books.Book, error) {
	return mockBooks, mockBooksErr
}

func (m *mockService) RemoveBooks(ids ...string) error {
	return mockBooksErr
}

func (m *mockService) UpdateBook(bk books.Book) error {
	return mockBooksErr
}
