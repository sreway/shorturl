// Package app configures and runs application.
package app

import (
	"context"
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

// Run shorturl application.
func Run(ctx context.Context) {
	var code int

	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("application", "shorturl")}))

	log.Info("start app")

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer func() {
		stop()
		if code != 0 {
			os.Exit(code)
		}
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

		configCache = cfg.GetStorage().GetCache()
		configPostgres = cfg.GetStorage().GetPostgres()
		configShortURL = cfg.GetShortURL()

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
			err = service.ProcQueue(ctx, cfg.GetShortURL().GetCheckTaskInterval())
			if err != nil {
				log.Error("failed processed task queue", err)
				stop()
				exit <- 1
				return
			}
		}()

		err = srv.Run(ctx, cfg.GetHTTP())
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
