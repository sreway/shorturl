package url

import (
	"net/url"

	"github.com/google/uuid"
)

//go:generate  mockgen -source=./internal/domain/url/url.go -destination=./internal/domain/url/mock/mock_url.go -package=urlMock
type (
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

func (e *entity) ID() uuid.UUID {
	return e.id
}

func (e *entity) LongURL() string {
	return e.longURL.String()
}

func (e *entity) ShortURL() string {
	return e.shortURL.String()
}

func (e *entity) LongValue() url.URL {
	return e.longURL
}

func (e *entity) ShortValue() url.URL {
	return e.shortURL
}

func (e *entity) UserID() uuid.UUID {
	return e.userID
}

func (e *entity) CorrelationID() string {
	return e.correlationID
}

func (e *entity) Deleted() bool {
	return e.deleted
}

func (e *entity) SetShortURL(value url.URL) {
	e.shortURL = value
}

func (e *entity) SetLongURL(value url.URL) {
	e.longURL = value
}

func (e *entity) SetCorrelationID(value string) {
	e.correlationID = value
}

func (e *entity) SetDeleted(value bool) {
	e.deleted = value
}

func NewURL(id, userID uuid.UUID) *entity {
	return &entity{
		id:     id,
		userID: userID,
	}
}
