package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	s.router.HandleFunc("/books", s.ListBooks)
	return &s
}

// ServeHTTP wraps router ServeHTTP calls calls.
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
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
	switch kind {
	case "request":
		code = http.StatusBadRequest
	}
	http.Error(w, fmt.Sprint(err), code)
}
