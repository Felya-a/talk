package errors

import "errors"

var (
	ErrValidation      = errors.New("data not valid")
	ErrCloseConnection = errors.New("close connection")
)
