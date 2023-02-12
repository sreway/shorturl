package http

import "errors"

var (
	ErrReadBody    = errors.New("read body error")
	ErrWriteBody   = errors.New("write body error")
	ErrInvalidSlug = errors.New("invalid slug")
)
