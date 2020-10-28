package database

import "errors"

type rdbError struct {
	message       string
	originalError error
}

func (e *rdbError) Error() string {
	return e.message + ":" + e.originalError.Error()
}

func (e *rdbError) Internal() bool {
	return true
}

var (
	ErrNotOpened = errors.New("error not opened db")
	ErrUnknown   = errors.New("error unknown")
)
