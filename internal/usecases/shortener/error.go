package shortener

import "errors"

var (
	// ErrParseURL implements shortener URL parsing error.
	ErrParseURL = errors.New("URL parsing error")
	// ErrDecodeURL implements shortener URL decoding error.
	ErrDecodeURL = errors.New("URL decoding error")
	// ErrParseUUID implements shortener UUID parsing error.
	ErrParseUUID = errors.New("UUID parsing error")
	// ErrTaskBufferFull implements shortener Utask buffer full error.
	ErrTaskBufferFull = errors.New("task buffer full")
)
