package url

import (
	"net/url"

	"github.com/google/uuid"
)

type (
	URL interface {
		ID() uuid.UUID
		UserID() uuid.UUID
		LongURL() string
		ShortURL() string
		LongValue() url.URL
		ShortValue() url.URL
		CorrelationID() string
		SetShortURL(value url.URL)
		SetCorrelationID(value string)
	}

	entity struct {
		id            uuid.UUID
		userID        uuid.UUID
		longURL       url.URL
		shortURL      url.URL
		correlationID string
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

func (e *entity) SetShortURL(value url.URL) {
	e.shortURL = value
}

func (e *entity) SetCorrelationID(value string) {
	e.correlationID = value
}

func NewURL(id, userID uuid.UUID, shortURL, longURL url.URL) *entity {
	return &entity{
		id:       id,
		shortURL: shortURL,
		longURL:  longURL,
		userID:   userID,
	}
}
