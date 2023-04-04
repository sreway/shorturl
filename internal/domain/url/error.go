package url

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrNotFound     = errors.New("URL not found")
	ErrAlreadyExist = errors.New("URL already exist")
	ErrDeleted      = errors.New("URL deleted")
)

type (
	ErrURL struct {
		error  error
		id     uuid.UUID
		userID uuid.UUID
	}
)

func (e *ErrURL) Error() string {
	return fmt.Sprintf("%s:%s", e.id.String(), e.error)
}

func (e ErrURL) Is(err error) bool {
	return errors.Is(e.error, err)
}

func (e *ErrURL) ID() uuid.UUID {
	return e.id
}

func (e *ErrURL) UserID() uuid.UUID {
	return e.userID
}

func NewURLErr(id, userID uuid.UUID, err error) error {
	return &ErrURL{
		id:     id,
		userID: userID,
		error:  err,
	}
}
