package authors

import (
	"bookshop/books"
	"encoding/json"
	"strings"
	"time"
)

const DateParsingFormat = "2006-01-02"

// Author is the model representing an author row in the datastore.
type Author struct {
	ID         string     `db:"id" json:"id,omitempty"`
	FirstName  string     `db:"first_name" json:"first_name,omitempty"`
	MiddleName string     `db:"middle_name" json:"middle_name,omitempty"`
	LastName   string     `db:"last_name" json:"last_name,omitempty"`
	DOB        *time.Time `db:"dob" json:"dob,omitempty"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at,omitempty"`

	Books []books.Book `json:"books,omitempty"`
}

// UnmarshalJSON is a custom unmarshaler that allows for passing in
// DOB (date of birth) values in a human readable format: "2006-01-02"
func (a *Author) UnmarshalJSON(b []byte) error {
	type value struct {
		DOB string `db:"dob" json:"dob,omitempty"`

		ID         string     `db:"id" json:"id,omitempty"`
		FirstName  string     `db:"first_name" json:"first_name,omitempty"`
		MiddleName string     `db:"middle_name" json:"middle_name,omitempty"`
		LastName   string     `db:"last_name" json:"last_name,omitempty"`
		UpdatedAt  *time.Time `db:"updated_at" json:"updated_at,omitempty"`

		Books []books.Book `json:"books,omitempty"`
	}
	var out value
	if err := json.Unmarshal(b, &out); err != nil {
		return err
	}

	dateFormat := DateParsingFormat

	// when data is coming directly from the database, it will include timezone
	// information so the format must handle this case
	if strings.Contains(out.DOB, "T00:00:00Z") {
		dateFormat = "2006-01-02T15:04:05Z"
	}

	var parsed *time.Time
	if out.DOB != "" {
		t, err := time.Parse(dateFormat, out.DOB)
		if err != nil {
			return err
		}
		parsed = &t
	}

	*a = Author{
		DOB: parsed,

		ID:         out.ID,
		FirstName:  out.FirstName,
		MiddleName: out.MiddleName,
		LastName:   out.LastName,
		Books:      out.Books,
	}
	return nil
}

// AuthorAndBook is the model representing a row of combined author and book data.
type AuthorAndBook struct {
	ID         string     `db:"id"`
	FirstName  string     `db:"first_name"`
	MiddleName string     `db:"middle_name"`
	LastName   string     `db:"last_name"`
	DOB        *time.Time `db:"dob"`
	UpdatedAt  *time.Time `db:"updated_at"`
	BookID     string     `db:"book_id"`
	BookTitle  string     `db:"book_title"`
	BookISBN   string     `db:"book_isbn"`
}

type BookAuth struct {
	BookID   string `db:"book_id" json:"book_id"`
	AuthorID string `db:"author_id" json:"author_id"`
}
