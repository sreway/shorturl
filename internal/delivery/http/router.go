package http

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/sreway/shorturl/docs"
	"github.com/sreway/shorturl/internal/config"
)

func (d *delivery) initRouter(http config.HTTP) *chi.Mux {
	var swaggerCfg config.Swagger = http.GetSwagger()
	docs.SwaggerInfo.Title = swaggerCfg.GetTitle()
	docs.SwaggerInfo.Description = swaggerCfg.GetTitle()
	docs.SwaggerInfo.Host = swaggerCfg.GetHost()
	docs.SwaggerInfo.BasePath = swaggerCfg.GetBasePath()
	docs.SwaggerInfo.Schemes = swaggerCfg.GetSchemes()

	router := chi.NewRouter()
	d.useMiddleware(http, router)
	d.routerURL(http, router)
	return router
}

func (d *delivery) routerURL(http config.HTTP, r chi.Router) {
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
		r.Route("/internal/stats", func(r chi.Router) {
			r.Use(trustedSubnet(http.GetTrustedSubnet()))
			r.Get("/", d.stats)
		})
	})

	r.Mount("/docs", httpSwagger.WrapHandler)
}
