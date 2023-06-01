package cache

import (
	"errors"
)

// ErrEmptyPath implements in-memory storage empty path error.
var ErrEmptyPath = errors.New("empty path")

// ErrInvalidStorageType implements in-memory storage invalid storage type error.
var ErrInvalidStorageType = errors.New("invalid storage type")
