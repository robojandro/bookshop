package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"bookshop/books"
	"bookshop/service"

	"github.com/gorilla/mux"
)

// HTTPServer is a struct for encompassing a book service and router for answering http requests.
type HTTPServer struct {
	svc    service.SVC
	router *mux.Router
}

// NewHTTPServer returns an HTTPServer type value and sets up routing.
func NewHTTPServer(svc service.SVC) *HTTPServer {
	r := mux.NewRouter()
	s := HTTPServer{
		svc:    svc,
		router: r,
	}

	bookRouter := s.router.PathPrefix("/books").Subrouter()
	{
		bookRouter.Methods(http.MethodGet).HandlerFunc(s.ListBooks)
		bookRouter.Methods(http.MethodPost).HandlerFunc(s.AddBook)
		bookRouter.Methods(http.MethodPatch).HandlerFunc(s.UpdateBook)

		bookRouter.Methods(http.MethodDelete).Path("/{book_id}").HandlerFunc(s.RemoveBook)
	}
	return &s
}

// ServeHTTP wraps router ServeHTTP calls calls.
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, r)
}

type bookBody struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	ISBN  string `json:"isbn"`
}

// AddBook adds a book generating its UUID if it doesn't already exist.
func (s *HTTPServer) AddBook(w http.ResponseWriter, r *http.Request) {
	var bk bookBody
	if err := json.NewDecoder(r.Body).Decode(&bk); err != nil {
		s.handleError(w, "request", err)
		return
	}

	created, err := s.svc.AddBook(bk.Title, bk.ISBN)
	if err != nil {
		s.handleError(w, "service", err)
		return
	}

	added, err := json.Marshal(created)
	if err != nil {
		s.handleError(w, "other", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	s.serve(w, added)
}

// ListBooks answers requests to list books.
func (s *HTTPServer) ListBooks(w http.ResponseWriter, r *http.Request) {
	bks, err := s.svc.ListBooks()
	if err != nil {
		s.handleError(w, "service", err)
		return
	}

	bkList, err := json.Marshal(bks)
	if err != nil {
		s.handleError(w, "other", err)
		return
	}

	s.serve(w, bkList)
}

// RemoveBook removes a book by its UUID if it exists.
func (s *HTTPServer) RemoveBook(w http.ResponseWriter, r *http.Request) {
	bkID := mux.Vars(r)["book_id"]
	if strings.TrimSpace(bkID) == "" {
		s.handleError(w, "request", errors.New("book ID cannot be blank"))
		return
	}

	if err := s.svc.RemoveBooks(bkID); err != nil {
		s.handleError(w, "service", err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	s.serve(w, []byte{})
}

// UpdateBook updates a book based on it's
func (s *HTTPServer) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var bk bookBody
	if err := json.NewDecoder(r.Body).Decode(&bk); err != nil {
		s.handleError(w, "request", err)
		return
	}

	if strings.TrimSpace(bk.ID) == "" {
		s.handleError(w, "request", errors.New("book ID cannot be blank"))
		return
	}

	err := s.svc.UpdateBook(books.Book{
		ID:    bk.ID,
		Title: bk.Title,
		ISBN:  books.ISBN(bk.ISBN),
	})
	if err != nil {
		s.handleError(w, "service", err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	s.serve(w, []byte{})
}

func (s *HTTPServer) serve(w http.ResponseWriter, v []byte) {
	_, err := w.Write(v)
	if err != nil {
		s.handleError(w, "other", err)
		return
	}
}

func (s *HTTPServer) handleError(w http.ResponseWriter, kind string, err error) {
	log.Print(err)
	code := http.StatusInternalServerError

	if kind == "service" && errors.Is(err, service.ErrDuplicate) {
		code = http.StatusUnprocessableEntity
	}

	if kind == "request" {
		code = http.StatusBadRequest
	}
	http.Error(w, fmt.Sprint(err), code)
}
