package apperror

import (
	"errors"
	"fmt"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
)

type ValidationError struct {
	Fields map[string]string
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("validation failed with %d errors", len(v.Fields))
}
