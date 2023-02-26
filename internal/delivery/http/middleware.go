package http

import (
	"compress/gzip"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sreway/shorturl/internal/config"
)

func (d *delivery) useMiddleware(http config.HTTP, r chi.Router) {
	r.Use(middleware.Compress(http.GetCompressLevel(), http.GetCompressTypes()...))
	r.Use(decodeGZIP)
}

func decodeGZIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				panic(err)
			}
			r.Body = reader
		}
		next.ServeHTTP(rw, r)
	})
}
