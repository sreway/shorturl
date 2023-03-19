package http

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"

	"github.com/sreway/shorturl/internal/config"
	"github.com/sreway/shorturl/internal/usecases"
)

type (
	delivery struct {
		shortener usecases.Shortener
		cfg       config.HTTP
		router    *chi.Mux
		logger    *slog.Logger
	}
)

func New(uc usecases.Shortener) *delivery {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "http")}))
	d := &delivery{
		shortener: uc,
		logger:    log,
	}
	return d
}

func (d *delivery) Run(ctx context.Context, config config.HTTP) error {
	d.router = d.initRouter(config)
	httpServer := &http.Server{
		Addr:    config.GetAddress(),
		Handler: d.router,
	}

	d.cfg = config

	ctxServer, stopServer := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		d.logger.Info("trigger graceful shutdown http server")
		err := httpServer.Shutdown(ctxServer)
		if err != nil {
			d.logger.Error("shutdown http server", err)
		}
		stopServer()
	}()
	d.logger.Info("http service is ready to listen and serv")
	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-ctxServer.Done()
	return nil
}
