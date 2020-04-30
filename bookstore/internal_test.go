package bookstore_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"bookstore/bookstore"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookstore(t *testing.T) {
	testDB, store := GetDB(t, "dev", "playground")
	LoadData(t, testDB, "./testdata/mockdb.json")
	defer cleanData(testDB, []string{"books"})
	t.Run("ReadBooks", func(t *testing.T) {
		books := store.ReadBooks()
		assert.NotNil(t, books)
		assert.Equal(t, "cb0b9721-7631-4b2a-94a2-493c559da893", books[0].ID)
		assert.Equal(t, "titleA", books[0].Title)
		assert.Equal(t, "9783161484100", string(books[0].ISBN))
		assert.Equal(t, "978-3-16-148410-0", books[0].ISBN.String())
	})
}

func GetDB(t *testing.T, dbUser, dbName string) (*sqlx.DB, bookstore.Store) {
	testDB, err := bookstore.ConnectDB(dbUser, dbName)
	require.NoError(t, err)
	return testDB, bookstore.NewBookstore(testDB)
}

type mockContents struct {
	Books []bookstore.Book `json:"books"`
}

func LoadData(t *testing.T, db *sqlx.DB, mockDataFile string) {
	data, err := ioutil.ReadFile(mockDataFile)
	require.NoError(t, err)
	//fmt.Printf("data: %s\n", data)

	var mocked mockContents
	err = json.Unmarshal(data, &mocked)
	require.NoError(t, err)

	//fmt.Printf("mocked: % #v\n", pretty.Formatter(mocked))
	err = insertBooks(db, mocked.Books)
	require.NoError(t, err)
}

func insertBooks(db *sqlx.DB, books []bookstore.Book) error {
	tx := db.MustBegin()
	bookIn :=
		`INSERT INTO books (id, title, isbn, created_at, updated_at)
		VALUES (:id, :title, :isbn, NOW(), NOW());`
	for _, book := range books {
		_, err := tx.NamedExec(bookIn, book)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func cleanData(db *sqlx.DB, tables []string) {
	tx := db.MustBegin()
	for _, table := range tables {
		log.Printf("clearing out table: %s\n", table)
		res := tx.MustExec(fmt.Sprintf("TRUNCATE TABLE %s", table))
		if res == nil {
			tx.Rollback()
			panic("failed truncating")
		}
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
