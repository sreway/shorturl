package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"

	"golang.org/x/exp/slog"

	"github.com/sreway/shorturl/internal/config"
	"github.com/sreway/shorturl/internal/delivery/http"
	"github.com/sreway/shorturl/internal/repository/storage/cache"
	"github.com/sreway/shorturl/internal/usecases/adapters/storage"
	"github.com/sreway/shorturl/internal/usecases/shortener"
)

func init() {
	var (
		httpServerAddress string
		shortenBaseURL    string
		storageCachePath  string
	)

	flag.StringVar(&httpServerAddress, "a", httpServerAddress,
		"http server address: scheme:host:port")
	flag.StringVar(&shortenBaseURL, "b", shortenBaseURL, "shorten base url")
	flag.StringVar(&storageCachePath, "f", storageCachePath, "storage cache file path")
	flag.Parse()

	_, exist := os.LookupEnv("BASE_URL")
	if !exist {
		_ = os.Setenv("BASE_URL", shortenBaseURL)
	}

	_, exist = os.LookupEnv("SERVER_ADDRESS")
	if !exist {
		_ = os.Setenv("SERVER_ADDRESS", httpServerAddress)
	}

	_, exist = os.LookupEnv("FILE_STORAGE_PATH")
	if !exist {
		_ = os.Setenv("FILE_STORAGE_PATH", storageCachePath)
	}
}

func main() {
	var code int

	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "main")}))

	log.Info("start app")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer func() {
		stop()
		os.Exit(code)
	}()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	exit := make(chan int)

	go func() {
		defer func() {
			wg.Done()
		}()

		var (
			cfg            config.Config
			configCache    config.Cache
			configShortURL config.ShortURL
			repo           storage.URL
		)

		cfg, err := config.NewConfig()
		if err != nil {
			log.Error("failed initialize config", err)
			stop()
			exit <- 1
			return
		}

		configCache = cfg.Storage().Cache()
		configShortURL = cfg.ShortURL()

		repo = cache.New(
			cache.Counter(configShortURL.GetCounter()),
			cache.File(configCache.GetFilePath()),
		)
		defer func() {
			err = repo.Close()
			if err != nil {
				log.Error("failed close url repository", err)
			}
		}()

		service := shortener.New(repo, cfg.ShortURL())
		srv := http.New(service)

		err = srv.Run(ctx, cfg.HTTP())
		if err != nil {
			log.Error("failed run delivery", err)
			stop()
			exit <- 1
			return
		}
	}()

	go func() {
		<-ctx.Done()
		stop()
		exit <- 0
		log.Info("trigger graceful shutdown app")
	}()

	code = <-exit
	wg.Wait()
}
