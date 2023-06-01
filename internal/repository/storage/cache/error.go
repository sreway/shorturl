package cache

import "errors"

var (
	// ErrEmptyPath implements in-memory storage empty path error.
	ErrEmptyPath = errors.New("empty path")
	// ErrInvalidStorageType implements in-memory storage invalid storage type error.
	ErrInvalidStorageType = errors.New("invalid storage type")
)
