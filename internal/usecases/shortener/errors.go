package shortener

import "errors"

var (
	ErrParseURL  = errors.New("URL parsing error")
	ErrDecodeURL = errors.New("URL decoding error")
)
