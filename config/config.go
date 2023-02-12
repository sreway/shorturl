package config

import (
	"net/url"

	"github.com/caarlos0/env/v6"
)

type (
	Config struct {
		Server Server
	}

	Server struct {
		HTTP     HTTP
		ShortURL ShortURL
	}

	HTTP struct {
		Scheme  string `env:"SERVER_SCHEME" envDefault:"http"`
		Address string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	}

	ShortURL struct {
		BaseURL *url.URL `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
		Counter uint64   `env:"COUNTER" envDefault:"1000000000"`
	}
)

func New() (*Config, error) {
	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
