package http

import (
	"errors"

	"github.com/go-chi/render"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrInternalServer = errors.New("internal server error")
	ErrStorageCheck   = errors.New("failed storage check")
)

// errRender implements renderer interface for managing response payloads.
func errRender(statusCode int, err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: statusCode,
		ErrorText:      err.Error(),
	}
}
