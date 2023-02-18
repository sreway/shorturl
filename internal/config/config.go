package config

import (
	"net/url"

	"github.com/caarlos0/env/v6"
)

type (
	Config interface {
		HTTP() HTTP
		ShortURL() ShortURL
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
		http     http
		shortURL shortURL
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

func (c *config) HTTP() HTTP {
	return &c.http
}

func (c *config) ShortURL() ShortURL {
	return &c.shortURL
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

func NewConfig() (Config, error) {
	cfg := new(config)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func NewHTTPConfig(scheme, address string) HTTP {
	return &http{
		scheme,
		address,
	}
}

func NewShortURLConfig(u *url.URL, counter uint64) ShortURL {
	return &shortURL{
		u,
		counter,
	}
}
