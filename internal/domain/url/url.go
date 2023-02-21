package url

import (
	"net/url"
)

type (
	URL interface {
		ID() uint64
		LongURL() *url.URL
		ShortURL() *url.URL
	}

	entity struct {
		id       uint64
		longURL  *url.URL
		shortURL *url.URL
	}
)

func (e *entity) ID() uint64 {
	return e.id
}

func (e *entity) LongURL() *url.URL {
	return e.longURL
}

func (e *entity) ShortURL() *url.URL {
	return e.shortURL
}

func NewURL(id uint64, shortURL, longURL *url.URL) *entity {
	return &entity{
		id:       id,
		shortURL: shortURL,
		longURL:  longURL,
	}
}
