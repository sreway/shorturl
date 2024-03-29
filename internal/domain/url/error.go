package url

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// ErrNotFound implements short URL not found error.
var ErrNotFound = errors.New("URL not found")

// ErrAlreadyExist implements short URL already exist error.
var ErrAlreadyExist = errors.New("URL already exist")

// ErrDeleted implements short URL already deleted error.
var ErrDeleted = errors.New("URL deleted")

// ErrURL defines short URL error.
type ErrURL struct {
	error  error
	id     uuid.UUID
	userID uuid.UUID
}

// Error implements "Error" method for short URL error.
func (e *ErrURL) Error() string {
	return fmt.Sprintf("%s:%s", e.id.String(), e.error)
}

// Is implements "Is" method for short URL error.
func (e ErrURL) Is(err error) bool {
	return errors.Is(e.error, err)
}

// ID implements getting short URL ID.
func (e *ErrURL) ID() uuid.UUID {
	return e.id
}

// UserID implements getting short URL user ID.
func (e *ErrURL) UserID() uuid.UUID {
	return e.userID
}

// NewURLErr implements the creation of the short URL error.
func NewURLErr(id, userID uuid.UUID, err error) error {
	return &ErrURL{
		id:     id,
		userID: userID,
		error:  err,
	}
}
