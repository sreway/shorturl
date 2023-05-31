package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"

	"golang.org/x/exp/slog"

	entity "github.com/sreway/shorturl/internal/domain/url"
	"github.com/sreway/shorturl/internal/usecases/shortener"
)

var urlSlug = regexp.MustCompile(`[^/][A-Za-z\d]+$`)

func (d *delivery) addURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	userID, ok := r.Context().Value(ctxKeyUserID{}).(string)
	if !ok {
		d.logger.Error("invalid user id", ErrInvalidRequest, slog.String("userID", userID))
		handelErrURL(w, ErrInvalidRequest)
		return
	}

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

	u, err := d.shortener.CreateURL(r.Context(), string(b), userID)
	if err != nil {
		handelErrURL(w, err)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	if err != nil && !errors.Is(err, entity.ErrAlreadyExist) {
		return
	}
	_, err = w.Write([]byte(u.ShortURL()))
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
	w.Header().Set("Location", u.LongURL())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (d *delivery) shortURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(ctxKeyUserID{}).(string)
	if !ok {
		d.logger.Error("invalid user id", ErrInvalidRequest, slog.String("userID", userID))
		handelErrURL(w, ErrInvalidRequest)
		return
	}

	req := new(shortURLRequest)
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		d.logger.Error("failed decode request url", err, slog.String("handler", "shortURL"))
		handelErrURL(w, ErrDecodeBody)
		return
	}

	res := new(shortURLResponse)

	u, err := d.shortener.CreateURL(r.Context(), req.URL, userID)
	if err != nil {
		handelErrURL(w, err)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	if err != nil && !errors.Is(err, entity.ErrAlreadyExist) {
		return
	}

	res.Result = u.ShortURL()

	// not use json encoder because it add new line for stream
	data, err := json.Marshal(res)
	if err != nil {
		d.logger.Error("failed marshal response url", err, slog.String("handler", "shortURL"))
		handelErrURL(w, err)
	}

	_, err = w.Write(data)
	if err != nil {
		d.logger.Error("write body", err, slog.String("handler", "shortURL"))
		handelErrURL(w, ErrWriteBody)
		return
	}
}

func (d *delivery) userURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(ctxKeyUserID{}).(string)
	if !ok {
		d.logger.Error("invalid user id", ErrInvalidRequest,
			slog.String("userID", userID), slog.String("handler", "getUserURLs"))
		handelErrURL(w, ErrInvalidRequest)
		return
	}

	urls, err := d.shortener.GetUserURLs(r.Context(), userID)
	if err != nil {
		d.logger.Error("failed get user urls", err,
			slog.String("userID", userID), slog.String("handler", "getUserURLs"))
		handelErrURL(w, err)
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make([]userURLResponse, 0, len(urls))

	for _, url := range urls {
		resp = append(resp, userURLResponse{
			url.ShortURL(),
			url.LongURL(),
		})
	}

	data, err := json.Marshal(resp)
	if err != nil {
		d.logger.Error("failed marshal response url", err, slog.String("handler", "getUserURLs"))
		handelErrURL(w, err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		d.logger.Error("write body", err, slog.String("handler", "getUserURLs"))
		handelErrURL(w, ErrWriteBody)
		return
	}
}

func (d *delivery) batchURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(ctxKeyUserID{}).(string)
	if !ok {
		d.logger.Error("invalid user id", ErrInvalidRequest,
			slog.String("userID", userID), slog.String("handler", "batchURL"))
		handelErrURL(w, ErrInvalidRequest)
		return
	}

	req := new([]batchURLRequest)
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		d.logger.Error("failed decode request url", err, slog.String("handler", "batchURL"))
		handelErrURL(w, ErrDecodeBody)
		return
	}

	correlationID := []string{}
	rawURL := []string{}

	for _, item := range *req {
		if item.CorrelationID != "" {
			correlationID = append(correlationID, item.CorrelationID)
		}

		if item.OriginalURL != "" {
			rawURL = append(rawURL, item.OriginalURL)
		}
	}

	if len(correlationID) != len(rawURL) {
		d.logger.Error("slice correlation id length is not equal to the length of raw slicer URLs",
			ErrInvalidRequest, slog.String("handler", "BatchURL"))
		handelErrURL(w, ErrInvalidRequest)
		return
	}

	urls, err := d.shortener.BatchURL(r.Context(), correlationID, rawURL, userID)
	if err != nil {
		d.logger.Error("failed batch add urls", err, slog.String("handler", "batchURL"))
		handelErrURL(w, err)
		return
	}

	resp := make([]batchURLResponse, 0)

	for _, i := range urls {
		resp = append(resp, batchURLResponse{
			i.CorrelationID(), i.ShortURL(),
		})
	}

	data, err := json.Marshal(resp)
	if err != nil {
		d.logger.Error("failed marshal response url", err, slog.String("handler", "batchURL"))
		handelErrURL(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		d.logger.Error("write body", err, slog.String("handler", "batchURL"))
		handelErrURL(w, ErrWriteBody)
		return
	}
}

func (d *delivery) deleteURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(ctxKeyUserID{}).(string)
	if !ok {
		d.logger.Error("invalid user id", ErrInvalidRequest,
			slog.String("userID", userID), slog.String("handler", "deleteBatchURL"))
		handelErrURL(w, ErrInvalidRequest)
		return
	}

	urls := new([]string)
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&urls); err != nil {
		d.logger.Error("failed decode request", err, slog.String("handler", "deleteURL"))
		handelErrURL(w, ErrDecodeBody)
		return
	}

	err := d.shortener.DeleteURL(r.Context(), userID, *urls)
	if err != nil {
		d.logger.Error("failed delete urls", err, slog.String("handler", "deleteURL"))
		handelErrURL(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (d *delivery) ping(w http.ResponseWriter, r *http.Request) {
	err := d.shortener.StorageCheck(r.Context())
	if err != nil {
		d.logger.Error("failed check storage", err, slog.String("handler", "ping"))
		handelErrURL(w, ErrStorageCheck)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	case errors.Is(err, shortener.ErrParseUUID):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, ErrInvalidRequest):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, entity.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, entity.ErrAlreadyExist):
		w.WriteHeader(http.StatusConflict)
	case errors.Is(err, ErrStorageCheck):
		w.WriteHeader(http.StatusInternalServerError)
	case errors.Is(err, entity.ErrDeleted):
		w.WriteHeader(http.StatusGone)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}
