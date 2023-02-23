package http

import "errors"

var (
	ErrReadBody    = errors.New("read body error")
	ErrEmptyBody   = errors.New("empty body")
	ErrWriteBody   = errors.New("write body error")
	ErrInvalidSlug = errors.New("invalid slug")
	ErrDecodeBody  = errors.New("failed decode body")
)
