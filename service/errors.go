package service

import (
	"errors"
	"fmt"
)

var (
	ErrDuplicate = NewErrDuplicate("")
)

func NewErrDuplicate(title string) error {
	return &DuplicateError{
		Err:   errors.New("isbn already exists with same title"),
		Title: title,
	}
}

type DuplicateError struct {
	Err   error
	Title string
}

func (e *DuplicateError) Error() string {
	return fmt.Sprintf("%s", e.Err) + ": " + e.Title
}

func (e *DuplicateError) Is(target error) bool {
	return target == ErrDuplicate
}
