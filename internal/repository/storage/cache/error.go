package cache

import "errors"

var (
	ErrEmptyPath          = errors.New("empty path")
	ErrInvalidStorageType = errors.New("invalid storage type")
)
