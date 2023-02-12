package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"

	"github.com/sreway/shorturl/config"
	"github.com/sreway/shorturl/internal/delivery/http"
	repo "github.com/sreway/shorturl/internal/repository/storage/cache/url"
	"github.com/sreway/shorturl/internal/usecases/shortener"
	log "github.com/sreway/shorturl/pkg/tools/logger"
)

func main() {
	var (
		err  error
		code int
	)
	log.Info("start app")

	cfg, err := config.New()
	if err != nil {
		log.Fatal("new config", zap.Error(err))
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
	service := shortener.New(repoURL, cfg.Server.ShortURL.BaseURL, cfg.Server.ShortURL.Counter)
	srv := http.New(service)

	go func() {
		defer wg.Done()
		err = srv.Run(ctx, &cfg.Server.HTTP)
		if err != nil {
			log.Error("delivery run", zap.Error(err))
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
				log.Error("trigger shutdown app")
				exit <- 1
			}
		}
	}()

	code = <-exit
}
