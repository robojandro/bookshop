package authors

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type AuthorStore struct {
	db *sqlx.DB
}

func NewAuthorStore(db *sqlx.DB) AuthorStore {
	return AuthorStore{db: db}
}

// ReadAuthors will return a list of all authors.
func (s *AuthorStore) ReadAuthors() ([]Author, error) {
	auths := []Author{}
	if err := s.db.Select(&auths, "SELECT * FROM authors ORDER BY last_name ASC"); err != nil {
		return nil, errors.Wrap(err, "failed to read authors")
	}
	return auths, nil
}

// ReadAuthor will return the given author looked up author ID.
func (s *AuthorStore) ReadAuthor(id string) (Author, error) {
	auth := Author{}
	row := s.db.QueryRowx("SELECT * FROM authors WHERE id=$1", id)
	if err := row.Scan(&auth.ID, &auth.FirstName, &auth.MiddleName, &auth.LastName, &auth.DOB, &auth.UpdatedAt); err != nil {
		return auth, errors.Wrap(err, "failed to read author")
	}
	return auth, nil
}

// DeleteAuthor will delete an author by their ID.
func (s *AuthorStore) DeleteAuthor(id string) error {
	if id == "" {
		return errors.New("no id submitted to delete")
	}
	tx := s.db.MustBegin()
	qry, args, err := sqlx.In("DELETE FROM authors WHERE id = ?", id)
	if err != nil {
		return errors.Wrap(err, "failed deleting authors")
	}
	rows, err := tx.Queryx(s.db.Rebind(qry), args...)
	if err != nil {
		return errors.Wrap(err, "failed deleting author")
	}
	rows.Close()
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// UpsertAuthors will modify or add the authors in the given list.
func (s *AuthorStore) UpsertAuthors(auths []Author) error {
	if len(auths) == 0 {
		return errors.New("no authors to upsert")
	}
	tx := s.db.MustBegin()
	const sqlSetPre = `INSERT INTO authors (id, first_name, middle_name, last_name, dob, updated_at) VALUES `

	const sqlSetPost = ` ON CONFLICT(id) DO
	UPDATE SET
		id = EXCLUDED.id,
		first_name = EXCLUDED.first_name,
		middle_name = EXCLUDED.middle_name,
		last_name = EXCLUDED.last_name,
		dob = EXCLUDED.dob,
		updated_at = EXCLUDED.updated_at
	RETURNING *;`
	const sqlValues = `(?,?,?,?,?,?)`

	var qryRows []string
	var qryArgs []interface{}
	for _, b := range auths {
		qryRows = append(qryRows, sqlValues)
		qryArgs = append(qryArgs, b.ID)
		qryArgs = append(qryArgs, b.FirstName)
		qryArgs = append(qryArgs, b.MiddleName)
		qryArgs = append(qryArgs, b.LastName)
		qryArgs = append(qryArgs, b.DOB)
		qryArgs = append(qryArgs, time.Now())
	}
	joinedRows := strings.Join(qryRows, ",")
	joinedQuery := sqlSetPre + joinedRows + sqlSetPost

	var err error
	var rows *sqlx.Rows
	if rows, err = tx.Queryx(s.db.Rebind(joinedQuery), qryArgs...); err != nil {
		return errors.Wrap(err, "failed upserting authors")
	}
	rows.Close()
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
