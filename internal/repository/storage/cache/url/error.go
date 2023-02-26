package url

import "errors"

var (
	ErrNotFound  = errors.New("not found")
	ErrEmptyPath = errors.New("empty path")
)
