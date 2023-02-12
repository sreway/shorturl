package http

import "github.com/go-chi/chi/v5"

func (d *Delivery) initRouter() *chi.Mux {
	router := chi.NewRouter()
	d.routerURL(router)
	return router
}

func (d *Delivery) routerURL(r chi.Router) {
	r.Post("/", d.AddURL)
	r.Get("/{id}", d.GetURL)
}
