package http

import (
	"compress/gzip"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/sreway/shorturl/internal/config"
	"github.com/sreway/shorturl/internal/delivery/http/cookies"
)

type ctxKeyUserID struct{}

func (d *delivery) useMiddleware(http config.HTTP, r chi.Router) {
	r.Use(middleware.Compress(http.GetCompressLevel(), http.GetCompressTypes()...))
	r.Use(decodeGZIP)
	r.Use(signCookie(http.GetCookie().SignID, http.GetCookie().SecretKey))
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

func signCookie(name string, secretKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			val, err := cookies.ReadSigned(r, name, secretKey)
			if err != nil {
				id := uuid.New()
				cookie := http.Cookie{
					Name:  name,
					Value: id.String(),
				}
				cookies.WriteSigned(w, cookie, secretKey)
				val = id.String()
			}
			ctx := context.WithValue(r.Context(), ctxKeyUserID{}, val)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
