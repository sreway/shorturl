// Package url implements and describes the type of short URL.
package url

import (
	"net/url"

	"github.com/google/uuid"
)

//go:generate  mockgen -source=./internal/domain/url/url.go -destination=./internal/domain/url/mock/mock_url.go -package=urlMock
type (
	// URL describes the implementation of the short URL type.
	URL interface {
		ID() uuid.UUID
		UserID() uuid.UUID
		LongURL() string
		ShortURL() string
		LongValue() url.URL
		ShortValue() url.URL
		CorrelationID() string
		Deleted() bool
		SetLongURL(value url.URL)
		SetShortURL(value url.URL)
		SetCorrelationID(value string)
		SetDeleted(value bool)
	}

	entity struct {
		id            uuid.UUID
		userID        uuid.UUID
		longURL       url.URL
		shortURL      url.URL
		correlationID string
		deleted       bool
	}
)

// ID implements getting short URL ID.
func (e *entity) ID() uuid.UUID {
	return e.id
}

// LongURL implements getting a long URL as string.
func (e *entity) LongURL() string {
	return e.longURL.String()
}

// ShortURL implements getting a short URL as a string.
func (e *entity) ShortURL() string {
	return e.shortURL.String()
}

// LongValue implements getting a long URL value.
func (e *entity) LongValue() url.URL {
	return e.longURL
}

// ShortValue implements getting a short URL value.
func (e *entity) ShortValue() url.URL {
	return e.shortURL
}

// UserID implements getting user ID.
func (e *entity) UserID() uuid.UUID {
	return e.userID
}

// CorrelationID implements getting correlation ID.
func (e *entity) CorrelationID() string {
	return e.correlationID
}

// Deleted implements getting the deletion attribute.
func (e *entity) Deleted() bool {
	return e.deleted
}

// SetShortURL implements the setting of a short URL value.
func (e *entity) SetShortURL(value url.URL) {
	e.shortURL = value
}

// SetLongURL implements the setting of a long URL value.
func (e *entity) SetLongURL(value url.URL) {
	e.longURL = value
}

// SetCorrelationID implements the setting of the correlation ID.
func (e *entity) SetCorrelationID(value string) {
	e.correlationID = value
}

// SetDeleted implements the setting of the deletion attribute.
func (e *entity) SetDeleted(value bool) {
	e.deleted = value
}

// NewURL implements the creation of the short URL type.
func NewURL(id, userID uuid.UUID) *entity {
	return &entity{
		id:     id,
		userID: userID,
	}
}
