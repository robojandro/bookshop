package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bookshop/authors"
	"bookshop/books"
	"bookshop/service"

	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPServer(t *testing.T) {
	t.Run("authors", func(t *testing.T) {
		//dt, err := time.Parse(authors.DateParsingFormat, "1970-01-01")
		//require.NoError(t, err)

		t.Run("GET", func(t *testing.T) {
			t.Run("happy", func(t *testing.T) {
				mockAuth = authors.Author{
					ID:         "auth01",
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
					//DOB:        &dt,
					Books: []books.Book{
						{
							ID:    "abc01",
							Title: "titleA",
							ISBN:  "9783161484100",
						},
					},
				}

				resp := makeRequest(t, "GET", "/authors/auth01", "")

				require.Equal(t, http.StatusOK, resp.Code)

				var content authors.Author
				fmt.Printf("resp.Body: %s\n", resp.Body)
				err := json.NewDecoder(resp.Body).Decode(&content)
				fmt.Printf("content: % #v\n", pretty.Formatter(content))
				require.NoError(t, err)
				assert.Equal(t, mockAuth.ID, content.ID)
				assert.Equal(t, mockAuth.FirstName, content.FirstName)
				assert.Equal(t, mockAuth.MiddleName, content.MiddleName)
				assert.Equal(t, mockAuth.LastName, content.LastName)

				fmt.Printf("content.Books: % #v\n", pretty.Formatter(content.Books))
				/*
					require.Len(t, content.Books, 1)
					assert.Equal(t, mockAuth.Books[0].Title, content.Books[0].Title)
				*/
			})
		})
	})

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

var mockAuth authors.Author
var mockAuthErr error

func (m *mockService) GetAuthor(id string) (authors.Author, error) {
	return mockAuth, mockAuthErr
}
