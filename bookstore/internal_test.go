// +build int

package bookstore_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"bookstore/bookstore"

	"github.com/jmoiron/sqlx"
	"github.com/robojandro/go-pgtesthelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookstore(t *testing.T) {
	h, dbh := initializeTestDB(t)
	defer h.CleanUp()

	store := bookstore.NewBookstore(dbh)

	t.Run("ReadBooks", func(t *testing.T) {
		books, err := store.ReadBooks()
		require.NoError(t, err)
		assert.NotNil(t, books)
		assert.Equal(t, "cb0b9721-7631-4b2a-94a2-493c559da893", books[0].ID)
		assert.Equal(t, "titleA", books[0].Title)
		assert.Equal(t, "9783161484100", string(books[0].ISBN))

		// test formatting of our to-string function
		assert.Equal(t, "978-3-16-148410-0", books[0].ISBN.String())
	})

	t.Run("InsertBooks", func(t *testing.T) {
		books := []bookstore.Book{
			{
				ID:    "abc02",
				Title: "titleB",
				ISBN:  "2222222222222",
			},
		}
		err := store.InsertBooks(books)
		require.NoError(t, err)

		results, err := store.ReadBooks()
		assert.Len(t, results, 2)
		assert.Equal(t, books[0].ID, results[1].ID)
		assert.Equal(t, books[0].Title, results[1].Title)
		assert.Equal(t, books[0].ISBN, results[1].ISBN)
	})

	t.Run("DeleteBooks", func(t *testing.T) {
		books := []bookstore.Book{
			{
				ID:    "def03",
				Title: "titleC",
				ISBN:  "3333333333333",
			},
			{
				ID:    "ghi04",
				Title: "titleD",
				ISBN:  "4444444444444",
			},
		}
		err := store.InsertBooks(books)
		require.NoError(t, err)

		err = store.DeleteBooks("def03", "ghi04")
		require.NoError(t, err)

		res, err := store.ReadBooks()
		require.NoError(t, err)

		found := map[string]bool{}
		for _, r := range res {
			found[r.ID] = true
		}
		assert.False(t, found["def03"])
		assert.False(t, found["ghi04"])
	})

	t.Run("Upsert", func(t *testing.T) {
		err := h.CleanTables([]string{"books"})
		require.NoError(t, err)

		books := []bookstore.Book{
			{
				ID:    "abc01",
				Title: "titleA",
				ISBN:  "1111111111111",
			},
		}
		err = store.InsertBooks(books)
		require.NoError(t, err)

		upsert := []bookstore.Book{
			{
				ID:    "abc01",
				Title: "titleXXXX",
				ISBN:  "9999999999999",
			},
			{
				ID:    "zzz00",
				Title: "titleZZZZ",
				ISBN:  "0000000000000",
			},
		}
		err = store.UpsertBooks(upsert)
		require.NoError(t, err)

		res, err := store.ReadBooks()
		require.NoError(t, err)
		require.Len(t, res, 2)

		assert.Equal(t, upsert[0].ID, res[0].ID)
		assert.Equal(t, upsert[0].Title, res[0].Title)
		assert.Equal(t, upsert[0].ISBN, res[0].ISBN)

		assert.Equal(t, upsert[1].ID, res[1].ID)
		assert.Equal(t, upsert[1].Title, res[1].Title)
		assert.Equal(t, upsert[1].ISBN, res[1].ISBN)
	})
}

func initializeTestDB(t *testing.T) (*pgtesthelper.Helper, *sqlx.DB) {
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

	dbh, err := h.CreateTempDB()

	mockData, err := ioutil.ReadFile("./testdata/mockdb.json")
	require.NoError(t, err)

	var data mockContents
	json.Unmarshal(mockData, &data)
	require.NoError(t, err)

	err = insertTestData(dbh, data)
	require.NoError(t, err)

	return &h, dbh
}

type mockContents struct {
	Books []bookstore.Book `json:"books"`
}

func insertTestData(dbh *sqlx.DB, data mockContents) error {
	tx := dbh.MustBegin()
	bookIn :=
		`INSERT INTO books (id, title, isbn, updated_at)
			        VALUES (:id, :title, :isbn, NOW());`
	for _, book := range data.Books {
		_, err := tx.NamedExec(bookIn, book)
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
