package shortener

import (
	"errors"
)

// ErrParseURL implements shortener URL parsing error.
var ErrParseURL = errors.New("URL parsing error")

// ErrDecodeURL implements shortener URL decoding error.
var ErrDecodeURL = errors.New("URL decoding error")

// ErrParseUUID implements shortener UUID parsing error.
var ErrParseUUID = errors.New("UUID parsing error")

// ErrTaskBufferFull implements shortener Utask buffer full error.
var ErrTaskBufferFull = errors.New("task buffer full")
