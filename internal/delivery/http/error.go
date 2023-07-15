package http

import (
	"errors"

	"github.com/go-chi/render"
)

// ErrInvalidRequest implements invalid http request error.
var ErrInvalidRequest = errors.New("invalid request")

// ErrInternalServer implements internal server error.
var ErrInternalServer = errors.New("internal server error")

// ErrStorageCheck implements storage check error.
var ErrStorageCheck = errors.New("failed storage check")

// ErrIPNotAllowed implements not allowed error.
var ErrIPNotAllowed = errors.New("ip not allowed")

// ErrTrustedSubnetNotSetup implements trusted subnet not setup error.
var ErrTrustedSubnetNotSetup = errors.New("trusted subnet not setup")

// ErrEmptyRealIPHeader implements missing X-Real-IP header.
var ErrEmptyRealIPHeader = errors.New("missing X-Real-IP header")

// errRender implements renderer interface for managing response payloads.
func errRender(statusCode int, err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: statusCode,
		ErrorText:      err.Error(),
	}
}
