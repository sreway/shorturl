package url

import "errors"

var (
	ErrParseURL = errors.New("failed parse url")
	ErrEmptyURL = errors.New("empty url data")
)
