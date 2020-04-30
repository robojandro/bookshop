package bookstore

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Store struct {
	db *sqlx.DB
}

func NewBookstore(db *sqlx.DB) Store {
	return Store{db: db}
}

func ConnectDB(user, dbname string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable", user, dbname))
	if err != nil {
		return nil, errors.Wrap(err, "db connection error")
	}
	return db, nil
}
