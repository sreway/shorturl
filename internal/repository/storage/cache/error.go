package cache

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrEmptyPath          = errors.New("empty path")
	ErrInvalidStorageType = errors.New("invalid storage type")
)
