package http

import (
	"io"
	"net/http"
	"regexp"

	"go.uber.org/zap"
)

var urlSlug = regexp.MustCompile(`[^\/][A-Za-z0-9]+$`)

func (d *Delivery) ShortURL(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		d.GetURL(w, r)
	case http.MethodPost:
		d.AddURL(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (d *Delivery) AddURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	b, err := io.ReadAll(r.Body)
	if err != nil {
		d.logger.Error("read body", zap.Error(err), zap.String("handler", "AddURL"))
		HandelErrURL(w, ErrReadBody)
		return
	}

	if len(b) == 0 {
		d.logger.Error("check len body", zap.Error(ErrEmptyBody), zap.String("handler", "AddURL"))
		HandelErrURL(w, ErrReadBody)
		return
	}

	u, err := d.shortener.CreateURL(r.Context(), string(b))
	if err != nil {
		HandelErrURL(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(u.ShortURL.String()))
	if err != nil {
		d.logger.Error("write body", zap.Error(err), zap.String("handler", "AddURL"))
		HandelErrURL(w, ErrWriteBody)
		return
	}
}

func (d *Delivery) GetURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if !urlSlug.Match([]byte(r.URL.Path)) {
		d.logger.Error("invalid slug", zap.Error(ErrInvalidSlug), zap.String("handler", "GetURL"))
		HandelErrURL(w, ErrInvalidSlug)
		return
	}

	id := urlSlug.Find([]byte(r.URL.Path))

	u, err := d.shortener.GetURL(r.Context(), string(id))
	if err != nil {
		HandelErrURL(w, err)
		return
	}
	w.Header().Set("Location", u.LongURL.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func HandelErrURL(w http.ResponseWriter, err error) {
	// temporarily do not check for an error, as a 400 code is always returned on error
	_ = err
	w.WriteHeader(http.StatusBadRequest)
}
