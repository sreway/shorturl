package shortener

import "errors"

var (
	ErrParseURL       = errors.New("URL parsing error")
	ErrDecodeURL      = errors.New("URL decoding error")
	ErrParseUUID      = errors.New("UUID parsing error")
	ErrTaskBufferFull = errors.New("task buffer full")
)
