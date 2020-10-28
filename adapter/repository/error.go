package repository

import "errors"

type repoError struct {
	message       string
	originalError error
}

func (e *repoError) Error() string {
	return e.message + ":" + e.originalError.Error()
}

func (e *repoError) Internal() bool {
	return true
}

var (
	ErrorStorage = errors.New("storage error")
	ErrUnknown   = errors.New("error unknown")
)
