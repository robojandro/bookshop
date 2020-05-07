package books

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type BookStore struct {
	db *sqlx.DB
}

func NewBookStore(db *sqlx.DB) BookStore {
	return BookStore{db: db}
}

// ReadBooks will return a list of all books.
func (s *BookStore) ReadBooks() ([]Book, error) {
	bks := []Book{}
	if err := s.db.Select(&bks, "SELECT * FROM books ORDER BY title ASC"); err != nil {
		return nil, errors.Wrap(err, "failed to read books")
	}
	return bks, nil
}

// ReadBookByISBN will return the given book looked up by isbn.
func (s *BookStore) ReadBookByISBN(isbn string) (Book, error) {
	bk := Book{}
	row := s.db.QueryRowx("SELECT * FROM books WHERE isbn=$1", isbn)
	if err := row.Scan(&bk.ID, &bk.Title, &bk.ISBN, &bk.UpdatedAt); err != nil {
		return bk, errors.Wrap(err, "failed to read book")
	}
	return bk, nil
}

// DeleteBooks will delete of books by their ID.
func (s *BookStore) DeleteBooks(ids ...string) error {
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
func (s *BookStore) UpsertBooks(bks []Book) error {
	if len(bks) == 0 {
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
	for _, b := range bks {
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
