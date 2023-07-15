package grpc

import (
	"errors"
)

// ErrInvalidRequest implements invalid grpc request error.
var ErrInvalidRequest = errors.New("invalid request")

// ErrInvalidUserID implements invalid user id error.
var ErrInvalidUserID = errors.New("invalid user id")

// ErrStorageCheck implements storage check error.
var ErrStorageCheck = errors.New("failed storage check")
