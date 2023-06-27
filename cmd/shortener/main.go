package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/sreway/shorturl/internal/app"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func buildInfo() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		buildVersion, buildDate, buildCommit)
}

func init() {
	var (
		httpServerAddress  string
		enableHTTPS        bool
		shortenBaseURL     string
		storageCachePath   string
		storagePostgresDSN string
		lookupEnv          = []string{
			"BASE_URL", "SERVER_ADDRESS", "FILE_STORAGE_PATH", "DATABASE_DSN", "ENABLE_HTTPS",
		}
	)

	buildInfo()
	flag.StringVar(&httpServerAddress, "a", httpServerAddress,
		"http server address: scheme:host:port")
	flag.StringVar(&shortenBaseURL, "b", shortenBaseURL, "shorten base url")
	flag.StringVar(&storageCachePath, "f", storageCachePath, "storage cache file path")
	flag.StringVar(&storagePostgresDSN, "d", storagePostgresDSN, "storage postgres dsn")
	flag.BoolVar(&enableHTTPS, "s", false, "enable https")
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
		case "ENABLE_HTTPS":
			_ = os.Setenv(env, strconv.FormatBool(enableHTTPS))
		}
	}
}

// @contact.name   API Support
// @contact.email  a.y.oleynik@gmail.com
func main() {
	ctx := context.Background()
	app.Run(ctx)
}
