package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/sreway/shorturl/internal/config"

	"golang.org/x/exp/slog"

	"github.com/sreway/shorturl/internal/delivery/http"
	repo "github.com/sreway/shorturl/internal/repository/storage/cache/url"
	"github.com/sreway/shorturl/internal/usecases/shortener"
)

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
		defer wg.Done()

		cfg, err := config.NewConfig()
		if err != nil {
			log.Error("failed initialize config", err)
			stop()
			exit <- 1
			return
		}
		repoURL := repo.New()
		service := shortener.New(repoURL, cfg.ShortURL())
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
