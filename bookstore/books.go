package bookstore

import (
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func (s *Store) ReadBooks() ([]Book, error) {
	books := []Book{}
	if err := s.db.Select(&books, "SELECT * FROM books ORDER BY title ASC"); err != nil {
		return nil, errors.Wrap(err, "failed to read books")
	}
	return books, nil
}
