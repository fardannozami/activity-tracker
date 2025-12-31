package errors

import "errors"

var (
	ErrClientEmailAlreadyExists = errors.New("email already registered")
)
