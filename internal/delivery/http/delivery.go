package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"go.uber.org/zap"

	"github.com/sreway/shorturl/config"
	"github.com/sreway/shorturl/internal/usecases"
	log "github.com/sreway/shorturl/pkg/tools/logger"
)

type (
	Delivery struct {
		shortener usecases.Shortener
		router    *chi.Mux
		logger    *zap.Logger
	}
)

func New(uc usecases.Shortener) *Delivery {
	l := log.GetLogger()
	d := &Delivery{
		shortener: uc,
		logger:    l.With(zap.String("service", "http")),
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
		log.Info("trigger graceful shutdown http server")
		err := httpServer.Shutdown(ctxServer)
		if err != nil {
			log.Fatal("shutdown http server", zap.Error(err))
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
