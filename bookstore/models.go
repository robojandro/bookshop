package bookstore

import (
	"fmt"
	"time"
)

type ISBN string

// String prints out the ISBN value in the proper dash delimited format.
//   example: 978-3-16-148410-0
func (i ISBN) String() string {
	digs := []byte(i)
	return fmt.Sprintf("%s-%s-%s-%s-%s", digs[0:3], digs[3:4], digs[4:6], digs[6:12], digs[len(digs)-1:])
}

// Book is the model representing an book row in the datastore.
type Book struct {
	ID        string     `db:"id" json:"id,omitempty"`
	Title     string     `db:"title" json:"title,omitempty"`
	ISBN      ISBN       `db:"isbn" json:"isbn,omitempty"`
	CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
