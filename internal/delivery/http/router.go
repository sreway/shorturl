package http

import "github.com/go-chi/chi/v5"

func (d *delivery) initRouter() *chi.Mux {
	router := chi.NewRouter()
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
