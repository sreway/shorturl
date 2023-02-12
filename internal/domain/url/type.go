package url

import (
	"net/url"
)

type URL struct {
	ID       uint64
	LongURL  *url.URL
	ShortURL *url.URL
}
