package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/exp/slog"

	"github.com/sreway/shorturl/config"
	"github.com/sreway/shorturl/internal/delivery/http"
	repo "github.com/sreway/shorturl/internal/repository/storage/cache/url"
	"github.com/sreway/shorturl/internal/usecases/shortener"
)

func main() {
	var (
		err  error
		code int
	)

	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "main")}))

	log.Info("start app")

	cfg, err := config.New()
	if err != nil {
		log.Error("failed initialize config", err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		wg.Wait()
		os.Exit(code)
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	exit := make(chan int)

	repoURL := repo.New()
	service := shortener.New(repoURL, &cfg.Server.ShortURL)
	srv := http.New(service)

	go func() {
		defer wg.Done()
		err = srv.Run(ctx, &cfg.Server.HTTP)
		if err != nil {
			log.Error("failed run delivery", err)
			signals <- syscall.SIGSTOP
		}
	}()

	go func() {
		for {
			s := <-signals
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Info("trigger graceful shutdown app")
				exit <- 0
			default:
				log.Info("trigger shutdown app")
				exit <- 1
			}
		}
	}()

	code = <-exit
}
