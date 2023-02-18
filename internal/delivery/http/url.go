package http

import (
	"io"
	"net/http"
	"regexp"

	"golang.org/x/exp/slog"
)

var urlSlug = regexp.MustCompile(`[^\/][A-Za-z0-9]+$`)

func (d *delivery) addURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	b, err := io.ReadAll(r.Body)
	if err != nil {
		d.logger.Error("read body", err, slog.String("handler", "AddURL"))
		HandelErrURL(w, ErrReadBody)
		return
	}

	if len(b) == 0 {
		d.logger.Error("check len body", ErrEmptyBody, slog.String("handler", "AddURL"))
		HandelErrURL(w, ErrReadBody)
		return
	}

	u, err := d.shortener.CreateURL(r.Context(), string(b))
	if err != nil {
		HandelErrURL(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(u.ShortURL().String()))
	if err != nil {
		d.logger.Error("write body", err, slog.String("handler", "AddURL"))
		HandelErrURL(w, ErrWriteBody)
		return
	}
}

func (d *delivery) getURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if !urlSlug.Match([]byte(r.URL.Path)) {
		d.logger.Error("invalid slug", ErrInvalidSlug, slog.String("handler", "GetURL"))
		HandelErrURL(w, ErrInvalidSlug)
		return
	}

	id := urlSlug.Find([]byte(r.URL.Path))

	u, err := d.shortener.GetURL(r.Context(), string(id))
	if err != nil {
		HandelErrURL(w, err)
		return
	}
	w.Header().Set("Location", u.LongURL().String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func HandelErrURL(w http.ResponseWriter, err error) {
	// temporarily do not check for an error, as a 400 code is always returned on error
	_ = err
	w.WriteHeader(http.StatusBadRequest)
}
