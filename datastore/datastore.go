package datastore

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Datastore struct {
	DB *sqlx.DB
}

func NewDatastore(db *sqlx.DB) Datastore {
	return Datastore{DB: db}
}

func ConnectDB(dbUser, dbPass, dbName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPass, dbName))
	if err != nil {
		return nil, errors.Wrap(err, "db connection error")
	}
	return db, nil
}
