package bookstore

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// ReadBooks will return a list of all books.
func (s *Store) ReadBooks() ([]Book, error) {
	books := []Book{}
	if err := s.db.Select(&books, "SELECT * FROM books ORDER BY title ASC"); err != nil {
		return nil, errors.Wrap(err, "failed to read books")
	}
	return books, nil
}

// InsertBooks will insert one book.
func (s *Store) InsertBooks(books []Book) error {
	tx := s.db.MustBegin()
	bookIns :=
		`INSERT INTO books (id, title, isbn, updated_at)
			        VALUES (:id, :title, :isbn, NOW());`
	for _, book := range books {
		_, err := tx.NamedExec(bookIns, book)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// DeleteBooks will delete of books by their ID.
func (s *Store) DeleteBooks(ids ...string) error {
	if len(ids) == 0 {
		return errors.New("no ids submitted to delete")
	}
	tx := s.db.MustBegin()
	qry, args, err := sqlx.In("DELETE FROM books WHERE id IN (?)", ids)
	if err != nil {
		return errors.Wrap(err, "failed deleting books")
	}
	rows, err := tx.Queryx(s.db.Rebind(qry), args...)
	if err != nil {
		return errors.Wrap(err, "failed deleting books")
	}
	rows.Close()
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// UpsertBooks will modify or add the books in the given list.
func (s *Store) UpsertBooks(books []Book) error {
	if len(books) == 0 {
		return errors.New("no books to upsert")
	}
	tx := s.db.MustBegin()
	const sqlSetPre = `INSERT INTO books (id, title, isbn, updated_at) VALUES `

	// updated_at = NOW() proves difficult to test, so manually set updated_at for updates
	const sqlSetPost = ` ON CONFLICT(id) DO
    UPDATE SET
        id = EXCLUDED.id,
        title = EXCLUDED.title,
        isbn = EXCLUDED.isbn,
        updated_at = EXCLUDED.updated_at
    RETURNING *;`
	const sqlValues = `(?,?,?,?)`

	var qryRows []string
	var qryArgs []interface{}
	for _, b := range books {
		qryRows = append(qryRows, sqlValues)
		qryArgs = append(qryArgs, b.ID)
		qryArgs = append(qryArgs, b.Title)
		qryArgs = append(qryArgs, b.ISBN)
		qryArgs = append(qryArgs, time.Now())
	}
	joinedRows := strings.Join(qryRows, ",")
	joinedQuery := sqlSetPre + joinedRows + sqlSetPost

	var err error
	var rows *sqlx.Rows
	if rows, err = tx.Queryx(s.db.Rebind(joinedQuery), qryArgs...); err != nil {
		return errors.Wrap(err, "failed upserting books")
	}
	rows.Close()
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
