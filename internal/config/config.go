package config

import (
	"net/url"

	"github.com/caarlos0/env/v6"
)

type (
	Config interface {
		HTTP() *http
		ShortURL() *shortURL
	}

	HTTP interface {
		GetScheme() string
		GetAddress() string
	}

	ShortURL interface {
		GetBaseURL() *url.URL
		GetCounter() uint64
	}

	config struct {
		http     *http
		shortURL *shortURL
	}

	http struct {
		Scheme  string `env:"SERVER_SCHEME" envDefault:"http"`
		Address string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	}

	shortURL struct {
		BaseURL *url.URL `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
		Counter uint64   `env:"COUNTER" envDefault:"1000000000"`
	}
)

func (c *config) HTTP() *http {
	return c.http
}

func (c *config) ShortURL() *shortURL {
	return c.shortURL
}

func (h *http) GetScheme() string {
	return h.Scheme
}

func (h *http) GetAddress() string {
	return h.Address
}

func (s *shortURL) GetBaseURL() *url.URL {
	return s.BaseURL
}

func (s *shortURL) GetCounter() uint64 {
	return s.Counter
}

func NewConfig() (*config, error) {
	cfgHTTP := new(http)
	cfgShortURL := new(shortURL)
	cfg := new(config)

	if err := env.Parse(cfgHTTP); err != nil {
		return nil, err
	}

	if err := env.Parse(cfgShortURL); err != nil {
		return nil, err
	}

	cfg.http = cfgHTTP
	cfg.shortURL = cfgShortURL
	return cfg, nil
}

func NewHTTPConfig(scheme, address string) *http {
	return &http{
		scheme,
		address,
	}
}

func NewShortURLConfig(u *url.URL, counter uint64) *shortURL {
	return &shortURL{
		u,
		counter,
	}
}
