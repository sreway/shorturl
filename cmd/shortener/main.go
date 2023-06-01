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
	"github.com/sreway/shorturl/internal/repository/storage/postgres"
	"github.com/sreway/shorturl/internal/usecases/adapters/storage"
	"github.com/sreway/shorturl/internal/usecases/shortener"
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
			configPostgres config.Postgres
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
		configPostgres = cfg.Storage().Postgres()
		configShortURL = cfg.ShortURL()

		switch {
		case len(configPostgres.GetDSN()) > 0:
			repo, err = postgres.New(ctx, configPostgres)
			if err == nil {
				log.Info("use postgres repository")
				break
			}
			log.Error("failed initialize postgres repository", err)
			fallthrough
		case len(configCache.GetFilePath()) > 0:
			repo = cache.New(cache.File(configCache.GetFilePath()))
			log.Info("use cache repository with specific file")
		default:
			repo = cache.New()
			log.Info("use default cache repository")
		}

		defer func() {
			err = repo.Close()
			if err != nil {
				log.Error("failed close url repository", err)
			}
		}()

		service := shortener.New(repo, configShortURL)
		srv := http.New(service)

		go func() {
			err = service.ProcQueue(ctx, cfg.ShortURL().GetCheckTaskInterval())
			if err != nil {
				log.Error("failed processed task queue", err)
				stop()
				exit <- 1
				return
			}
		}()

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
