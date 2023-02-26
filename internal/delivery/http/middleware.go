package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sreway/shorturl/internal/config"
)

func (d *delivery) useCompress(http config.HTTP, r chi.Router) {
	r.Use(middleware.Compress(http.GetCompressLevel(), http.GetCompressTypes()...))
}
