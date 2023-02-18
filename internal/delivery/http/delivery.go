package http

import (
	"context"
	"errors"
	"net/http"
	"os"

	"golang.org/x/exp/slog"

	"github.com/go-chi/chi/v5"

	"github.com/sreway/shorturl/config"
	"github.com/sreway/shorturl/internal/usecases"
)

type (
	Delivery struct {
		shortener usecases.Shortener
		router    *chi.Mux
		logger    *slog.Logger
	}
)

func New(uc usecases.Shortener) *Delivery {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "http")}))
	d := &Delivery{
		shortener: uc,
		logger:    log,
	}
	d.router = d.initRouter()
	return d
}

func (d *Delivery) Run(ctx context.Context, config *config.HTTP) error {
	httpServer := &http.Server{
		Addr:    config.Address,
		Handler: d.router,
	}

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
