// +build int

package bookstore_test

import (
	"encoding/json"
	"os"
	"testing"

	"bookstore/bookstore"

	"github.com/jmoiron/sqlx"
	"github.com/robojandro/go-pgtesthelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookstore(t *testing.T) {
	h := initializeTestDB(t)
	defer h.CleanUp()

	store := bookstore.NewBookstore(h.TestDB())

	t.Run("ReadBooks", func(t *testing.T) {
		books, err := store.ReadBooks()
		require.NoError(t, err)
		assert.NotNil(t, books)
		assert.Equal(t, "cb0b9721-7631-4b2a-94a2-493c559da893", books[0].ID)
		assert.Equal(t, "titleA", books[0].Title)
		assert.Equal(t, "9783161484100", string(books[0].ISBN))
		assert.Equal(t, "978-3-16-148410-0", books[0].ISBN.String())
	})
}

func initializeTestDB(t *testing.T) *pgtesthelper.Helper {
	var (
		schemaPath = "../sql/authors_books.sql"
		keepDB     = false
		dbPrefix   = "bookstore_testing"
		dbUser     = ""
		dbPass     = ""
	)

	if dbUser = os.Getenv("bookstore_dbuser"); dbUser == "" {
		t.Skip("missing env variable bookstore_dbuser")
	}
	if dbPass = os.Getenv("bookstore_dbpass"); dbPass == "" {
		t.Skip("missing env variable bookstore_dbuser")
	}

	h, err := pgtesthelper.NewHelper(schemaPath, dbPrefix, dbUser, dbPass, keepDB)
	require.NoError(t, err)

	err = h.CreateTestingDB()
	require.NoError(t, err)

	mockDB := "./testdata/mockdb.json"
	err = h.ParseMockData(mockDB, func(mockData []byte) error {
		return json.Unmarshal(mockData, &data)
	})
	require.NoError(t, err)

	err = h.LoadData("./testdata/mockdb.json", insertTestData)
	require.NoError(t, err)

	return &h
}

type mockContents struct {
	Books []bookstore.Book `json:"books"`
}

var data mockContents

var insertTestData = func(db *sqlx.DB) error {
	tx := db.MustBegin()
	bookIn :=
		`INSERT INTO books (id, title, isbn, created_at, updated_at)
			        VALUES (:id, :title, :isbn, NOW(), NOW());`
	for _, book := range data.Books {
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
