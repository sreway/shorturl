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
		r.Get("/ping", d.ping)
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/shorten", func(r chi.Router) {
			r.Post("/", d.shortURL)
			r.Post("/batch", d.batchURL)
		})
		r.Route("/user", func(r chi.Router) {
			r.Get("/urls", d.userURL)
			r.Delete("/urls", d.deleteURL)
		})
	})
}
