package main

import (
	"context"
	"flag"
	"os"

	"github.com/sreway/shorturl/internal/app"
)

func init() {
	var (
		httpServerAddress  string
		shortenBaseURL     string
		storageCachePath   string
		storagePostgresDSN string
		lookupEnv          = []string{
			"BASE_URL", "SERVER_ADDRESS", "FILE_STORAGE_PATH", "DATABASE_DSN",
		}
	)

	flag.StringVar(&httpServerAddress, "a", httpServerAddress,
		"http server address: scheme:host:port")
	flag.StringVar(&shortenBaseURL, "b", shortenBaseURL, "shorten base url")
	flag.StringVar(&storageCachePath, "f", storageCachePath, "storage cache file path")
	flag.StringVar(&storagePostgresDSN, "d", storagePostgresDSN, "storage postgres dsn")
	flag.Parse()

	for _, env := range lookupEnv {
		_, exist := os.LookupEnv(env)
		if exist {
			continue
		}

		switch env {
		case "BASE_URL":
			_ = os.Setenv(env, shortenBaseURL)
		case "SERVER_ADDRESS":
			_ = os.Setenv(env, httpServerAddress)
		case "FILE_STORAGE_PATH":
			_ = os.Setenv(env, storageCachePath)
		case "DATABASE_DSN":
			_ = os.Setenv(env, storagePostgresDSN)
		}
	}
}

// @contact.name   API Support
// @contact.email  a.y.oleynik@gmail.com
func main() {
	ctx := context.Background()
	app.Run(ctx)
}
