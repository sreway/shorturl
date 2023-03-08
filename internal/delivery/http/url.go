package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"

	"golang.org/x/exp/slog"

	shortURL "github.com/sreway/shorturl/internal/delivery/http/url"
	"github.com/sreway/shorturl/internal/usecases/shortener"
)

var urlSlug = regexp.MustCompile(`[^\/][A-Za-z0-9]+$`)

func (d *delivery) addURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	b, err := io.ReadAll(r.Body)
	if err != nil {
		d.logger.Error("read body", err, slog.String("handler", "AddURL"))
		handelErrURL(w, ErrReadBody)
		return
	}

	if len(b) == 0 {
		d.logger.Error("check len body", ErrEmptyBody, slog.String("handler", "AddURL"))
		handelErrURL(w, ErrEmptyBody)
		return
	}

	u, err := d.shortener.CreateURL(r.Context(), string(b))
	if err != nil {
		handelErrURL(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(u.ShortURL().String()))
	if err != nil {
		d.logger.Error("write body", err, slog.String("handler", "AddURL"))
		handelErrURL(w, ErrWriteBody)
		return
	}
}

func (d *delivery) getURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if !urlSlug.Match([]byte(r.URL.Path)) {
		d.logger.Error("invalid slug", ErrInvalidSlug, slog.String("handler", "GetURL"))
		handelErrURL(w, ErrInvalidSlug)
		return
	}

	id := urlSlug.Find([]byte(r.URL.Path))

	u, err := d.shortener.GetURL(r.Context(), string(id))
	if err != nil {
		handelErrURL(w, err)
		return
	}
	w.Header().Set("Location", u.LongURL().String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (d *delivery) shortURL(w http.ResponseWriter, r *http.Request) {
	var reqURL shortURL.Request
	w.Header().Set("Content-Type", "application/json")
	reqURL = shortURL.NewURLRequest(nil)
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&reqURL); err != nil {
		d.logger.Error("failed decode request url", err, slog.String("handler", "shortURL"))
		handelErrURL(w, ErrDecodeBody)
		return
	}

	u, err := d.shortener.CreateURL(r.Context(), reqURL.URL().String())
	if err != nil {
		handelErrURL(w, err)
		return
	}

	respURL := shortURL.NewURLResponse(u.ShortURL())

	// not use json encoder because it add new line for stream
	data, err := json.Marshal(respURL)
	if err != nil {
		d.logger.Error("failed marshal response url", err, slog.String("handler", "shortURL"))
		handelErrURL(w, err)
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		d.logger.Error("write body", err, slog.String("handler", "shortURL"))
		handelErrURL(w, ErrWriteBody)
		return
	}
}

func handelErrURL(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrEmptyBody):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, ErrInvalidSlug):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, ErrWriteBody):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, ErrReadBody):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, ErrDecodeBody):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, shortener.ErrDecodeURL):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, shortener.ErrParseURL):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, shortURL.ErrParseURL):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, shortURL.ErrEmptyURL):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}
