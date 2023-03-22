package url

import (
	"net/url"
)

type (
	URL interface {
		ID() [16]byte
		UserID() [16]byte
		LongURL() *url.URL
		ShortURL() *url.URL
		CorrelationID() string
		SetShortURL(value *url.URL)
		SetCorrelationID(value string)
	}

	entity struct {
		id            [16]byte
		userID        [16]byte
		longURL       *url.URL
		shortURL      *url.URL
		correlationID string
	}
)

func (e *entity) ID() [16]byte {
	return e.id
}

func (e *entity) LongURL() *url.URL {
	return e.longURL
}

func (e *entity) ShortURL() *url.URL {
	return e.shortURL
}

func (e *entity) UserID() [16]byte {
	return e.userID
}

func (e *entity) CorrelationID() string {
	return e.correlationID
}

func (e *entity) SetShortURL(value *url.URL) {
	e.shortURL = value
}

func (e *entity) SetCorrelationID(value string) {
	e.correlationID = value
}

func NewURL(id, userID [16]byte, shortURL, longURL *url.URL) *entity {
	return &entity{
		id:       id,
		shortURL: shortURL,
		longURL:  longURL,
		userID:   userID,
	}
}
