// +build int

package authors_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"bookshop/authors"
	"bookshop/books"

	"github.com/jmoiron/sqlx"
	"github.com/robojandro/go-pgtesthelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthors(t *testing.T) {
	h, dbh := initializeTestDB(t)
	defer h.CleanUp()

	store := authors.NewAuthorStore(dbh)

	dt, err := time.Parse(authors.DateParsingFormat, "1970-01-01")
	require.NoError(t, err)

	t.Run("ReadAuthors", func(t *testing.T) {
		auths, err := store.ReadAuthors()
		require.NoError(t, err)
		assert.NotNil(t, auths)
		assert.Equal(t, "0b5babb0-96d8-11ea-bb37-0242ac130002", auths[0].ID)
		assert.Equal(t, "first", auths[0].FirstName)
		assert.Equal(t, "1970-01-01", auths[0].DOB.Format(authors.DateParsingFormat))
	})

	t.Run("ReadAuthorAndBooks", func(t *testing.T) {
		auth, err := store.ReadAuthorAndBooks("0b5babb0-96d8-11ea-bb37-0242ac130002")
		require.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, "0b5babb0-96d8-11ea-bb37-0242ac130002", auth.ID)
		assert.Equal(t, "first", auth.FirstName)
		assert.Equal(t, "1970-01-01", auth.DOB.Format(authors.DateParsingFormat))

		require.Len(t, auth.Books, 2)
		assert.Equal(t, "titleA", auth.Books[0].Title)
		assert.Equal(t, "978-3-16-148410-0", auth.Books[0].ISBN.String())
		assert.Equal(t, "titleB", auth.Books[1].Title)
		assert.Equal(t, "444-4-44-444444-4", auth.Books[1].ISBN.String())
	})

	t.Run("DeleteBooks", func(t *testing.T) {
		upsert := []authors.Author{
			{
				ID:         "def03",
				FirstName:  "first3",
				MiddleName: "middle3",
				LastName:   "last3",
				DOB:        &dt,
			},
		}
		err = store.UpsertAuthors(upsert)
		require.NoError(t, err)

		err = store.DeleteAuthor("def03")
		require.NoError(t, err)

		res, err := store.ReadAuthors()
		require.NoError(t, err)

		found := map[string]bool{}
		for _, r := range res {
			found[r.ID] = true
		}
		assert.False(t, found["def03"])
	})

	t.Run("Upsert", func(t *testing.T) {
		err := h.CleanTables([]string{"books", "authors"})
		require.NoError(t, err)

		dt, err := time.Parse(authors.DateParsingFormat, "1970-01-01")
		require.NoError(t, err)

		auths := mockContents{
			Authors: []authors.Author{
				{
					ID:         "abc01",
					FirstName:  "first",
					MiddleName: "middle",
					LastName:   "last",
					DOB:        &dt,
				},
			},
		}

		err = insertTestData(dbh, auths)
		require.NoError(t, err)

		upsert := []authors.Author{
			{
				ID:         "abc01",
				FirstName:  "update",
				MiddleName: "update",
				LastName:   "update",
				DOB:        &dt,
			},
			{
				ID:         "zzz00",
				FirstName:  "new",
				MiddleName: "new",
				LastName:   "new",
				DOB:        &dt,
			},
		}
		err = store.UpsertAuthors(upsert)
		require.NoError(t, err)

		res, err := store.ReadAuthors()
		require.NoError(t, err)
		require.Len(t, res, 2)

		for i, u := range upsert {
			assert.Equal(t, upsert[i].ID, u.ID)
			assert.Equal(t, upsert[i].FirstName, u.FirstName)
			assert.Equal(t, upsert[i].MiddleName, u.MiddleName)
			assert.Equal(t, upsert[i].LastName, u.LastName)
			assert.Equal(t, upsert[i].DOB, u.DOB)
		}
	})
}

func initializeTestDB(t *testing.T) (*pgtesthelper.Helper, *sqlx.DB) {
	var (
		schemaPath = "../sql/authors_books.sql"
		keepDB     = false
		dbPrefix   = "books_testing"
		dbUser     = ""
		dbPass     = ""
	)

	if dbUser = os.Getenv("bookshop_dbuser"); dbUser == "" {
		t.Skip("missing env variable bookshop_dbuser")
	}
	if dbPass = os.Getenv("bookshop_dbpass"); dbPass == "" {
		t.Skip("missing env variable bookshop_dbpass")
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
	Authors   []authors.Author   `json:"authors"`
	BookAuths []authors.BookAuth `json:"books_authors"`
	Books     []books.Book       `json:"books"`
}

func insertTestData(dbh *sqlx.DB, data mockContents) error {
	tx := dbh.MustBegin()
	authIn :=
		`INSERT INTO authors (id, first_name, middle_name, last_name, dob, updated_at)
				  VALUES (:id, :first_name, :middle_name, :last_name, :dob, NOW());`
	for _, auth := range data.Authors {
		_, err := tx.NamedExec(authIn, auth)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

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

	time.Sleep(400 * time.Millisecond)
	tx2 := dbh.MustBegin()
	bookAuthIn :=
		`INSERT INTO books_authors (book_id, author_id) VALUES (:book_id, :author_id);`
	for _, bkAuth := range data.BookAuths {
		_, err := tx2.NamedExec(bookAuthIn, bkAuth)
		if err != nil {
			if err := tx2.Rollback(); err != nil {
				return err
			}
			return err
		}
	}
	if err := tx2.Commit(); err != nil {
		return err
	}
	return nil
}
