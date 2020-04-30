package bookstore

import (
	_ "github.com/lib/pq"
)

func (s *Store) ReadBooks() []Book {
	books := []Book{}
	s.db.Select(&books, "SELECT * FROM books ORDER BY title ASC")
	return books
}
