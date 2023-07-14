package http

import (
	"compress/gzip"
	"context"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"github.com/sreway/shorturl/internal/config"
	"github.com/sreway/shorturl/internal/delivery/http/cookies"
)

// ctxKeyUserID describes the type context value of the user ID.
type ctxKeyUserID struct{}

// useMiddleware implements middleware connection.
func (d *delivery) useMiddleware(http config.HTTP, r chi.Router) {
	r.Use(middleware.Compress(http.GetCompressLevel(), http.GetCompressTypes()...))
	r.Use(decodeGZIP)
	r.Use(signCookie(http.GetCookie().SignID, http.GetCookie().SecretKey))
}

// decodeGZIP implements compression middleware.
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

// signCookie implements sign cookie middleware.
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

// trustedSubnet implements validate trusted subnet middleware.
func trustedSubnet(subnet *net.IPNet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var emptySubnet *net.IPNet

			if subnet == emptySubnet {
				err := render.Render(w, r, errRender(http.StatusForbidden, ErrTrustedSubnetNotSetup))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				return
			}
			rip := r.Header.Get("X-Real-IP")
			if len(rip) == 0 {
				err := render.Render(w, r, errRender(http.StatusForbidden, ErrEmptyRealIPHeader))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				return
			}
			ip := net.ParseIP(rip)
			if !subnet.Contains(ip) {
				err := render.Render(w, r, errRender(http.StatusForbidden, ErrIPNotAllowed))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
