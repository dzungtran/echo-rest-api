package constants

import "errors"

var (
	ErrNotFound     error = errors.New("not found")
	ErrDuplicated   error = errors.New("duplicated")
	ErrUnauthorized error = errors.New("unauthorized")
	ErrForbidden    error = errors.New("forbidden")
)
