// Package http implements and describes the http server of the application.
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
		router    *chi.Mux
		logger    *slog.Logger
	}
)

// New implements http server initialization.
func New(uc usecases.Shortener) *delivery {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "http")}))
	d := &delivery{
		shortener: uc,
		logger:    log,
	}
	return d
}

// Run implements run http server.
func (d *delivery) Run(ctx context.Context, config config.HTTP) error {
	var err error
	d.router = d.initRouter(config)
	httpServer := &http.Server{
		Addr:    config.GetAddress(),
		Handler: d.router,
	}

	ctxServer, stopServer := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		d.logger.Info("trigger graceful shutdown http server")
		err = httpServer.Shutdown(ctxServer)
		if err != nil {
			d.logger.Error("shutdown http server", err)
		}
		stopServer()
	}()
	d.logger.Info("http service is ready to listen and serv")

	if config.GetScheme() == "https" {
		tlsCfg := config.GetTLS()
		err = httpServer.ListenAndServeTLS(tlsCfg.CertPath, tlsCfg.KeyPath)
	} else {
		err = httpServer.ListenAndServe()
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-ctxServer.Done()
	return nil
}
