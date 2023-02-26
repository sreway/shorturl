package http

import (
	"github.com/go-chi/chi/v5"

	"github.com/sreway/shorturl/internal/config"
)

func (d *delivery) initRouter(http config.HTTP) *chi.Mux {
	router := chi.NewRouter()
	d.useMiddleware(http, router)
	d.routerURL(router)
	return router
}

func (d *delivery) routerURL(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Post("/", d.addURL)
		r.Get("/{id}", d.getURL)
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", d.shortURL)
	})
}
