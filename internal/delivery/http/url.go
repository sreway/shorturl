package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"

	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	entity "github.com/sreway/shorturl/internal/domain/url"
	"github.com/sreway/shorturl/internal/usecases/shortener"
)

var urlSlug = regexp.MustCompile(`[^/][A-Za-z\d]+$`)

// addURL godoc
// @Summary add short URL
// @Description add short URL
// @ID addURL
// @Produce text/plain
// @Param longURL body string true "long URL to shorten"
// @Success 201 {string} string
// @Failure 409 {object} errResponse
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Failure 501 {object} errResponse
// @Router / [post]
func (d *delivery) addURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	userID, ok := r.Context().Value(ctxKeyUserID{}).(string)
	if !ok {
		d.logger.Error("invalid user id", ErrInvalidRequest, slog.String("userID", userID))
		d.handelErrURL(w, r, ErrInvalidRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		d.logger.Error("read body", err, slog.String("handler", "AddURL"))
		d.handelErrURL(w, r, ErrInvalidRequest)
		return
	}

	if len(b) == 0 {
		d.logger.Error("check len body", err, slog.String("handler", "AddURL"))
		d.handelErrURL(w, r, ErrInvalidRequest)
		return
	}

	u, err := d.shortener.CreateURL(r.Context(), string(b), userID)
	if err != nil {
		d.handelErrURL(w, r, err)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	if err != nil && !errors.Is(err, entity.ErrAlreadyExist) {
		return
	}
	_, err = w.Write([]byte(u.ShortURL()))
	if err != nil {
		d.logger.Error("write body", err, slog.String("handler", "AddURL"))
		d.handelErrURL(w, r, ErrInternalServer)
		return
	}
}

// getURL godoc
// @Summary get short URL
// @Description get short URL
// @ID getURL
// @Produce text/plain
// @Param id path string true "short URL id"
// @Success 200 {string} string
// @Failure 404 {object} errResponse
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Failure 501 {object} errResponse
// @Router /{id} [post]
func (d *delivery) getURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if !urlSlug.Match([]byte(r.URL.Path)) {
		d.logger.Error("invalid slug", ErrInvalidRequest, slog.String("handler", "getURL"))
		d.handelErrURL(w, r, ErrInvalidRequest)
		return
	}

	id := urlSlug.Find([]byte(r.URL.Path))

	u, err := d.shortener.GetURL(r.Context(), string(id))
	if err != nil {
		d.handelErrURL(w, r, err)
		return
	}
	w.Header().Set("Location", u.LongURL())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// shortURL godoc
// @Summary create short URL
// @Description create short URL
// @ID shortURL
// @Produce application/json
// @Param longURL body shortURLRequest true "long URL to shorten"
// @Success 201 {object} shortURLResponse
// @Failure 409 {object} errResponse
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Failure 501 {object} errResponse
// @Router /api/shorten [post]
func (d *delivery) shortURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	userID, ok := r.Context().Value(ctxKeyUserID{}).(string)
	if !ok {
		d.logger.Error("invalid user id", ErrInvalidRequest, slog.String("userID", userID))
		d.handelErrURL(w, r, ErrInvalidRequest)
		return
	}

	req := new(shortURLRequest)
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		d.logger.Error("failed decode request url", err, slog.String("handler", "shortURL"))
		d.handelErrURL(w, r, ErrInvalidRequest)
		return
	}

	res := new(shortURLResponse)

	u, err := d.shortener.CreateURL(r.Context(), req.URL, userID)
	if err != nil {
		d.handelErrURL(w, r, err)
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
		d.handelErrURL(w, r, err)
	}

	_, err = w.Write(data)
	if err != nil {
		d.logger.Error("write body", err, slog.String("handler", "shortURL"))
		d.handelErrURL(w, r, ErrInternalServer)
		return
	}
}

// userURL godoc
// @Summary get short URLs for user ID
// @Description get short URLs for user ID
// @ID userURL
// @Produce application/json
// @Success 201 {object} []userURLResponse
// @Failure 404 {object} errResponse
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Failure 501 {object} errResponse
// @Router /api/shorten/user/urls [get]
func (d *delivery) userURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	userID, ok := r.Context().Value(ctxKeyUserID{}).(string)
	if !ok {
		d.logger.Error("invalid user id", ErrInvalidRequest,
			slog.String("userID", userID), slog.String("handler", "getUserURLs"))
		d.handelErrURL(w, r, ErrInvalidRequest)
		return
	}

	urls, err := d.shortener.GetUserURLs(r.Context(), userID)
	if err != nil {
		d.logger.Error("failed get user urls", err,
			slog.String("userID", userID), slog.String("handler", "getUserURLs"))
		d.handelErrURL(w, r, err)
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make([]userURLResponse, len(urls))

	for idx, url := range urls {
		resp[idx] = userURLResponse{
			url.ShortURL(),
			url.LongURL(),
		}
	}

	data, err := json.Marshal(resp)
	if err != nil {
		d.logger.Error("failed marshal response url", err, slog.String("handler", "getUserURLs"))
		d.handelErrURL(w, r, err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		d.logger.Error("write body", err, slog.String("handler", "getUserURLs"))
		d.handelErrURL(w, r, ErrInternalServer)
		return
	}
}

// batchURL godoc
// @Summary create of several short URLs
// @Description create of several short URLs
// @ID batchURL
// @Produce application/json
// @Param bathURL body []batchURLRequest true "several long URLs to shorten"
// @Success 201 {object} []batchURLResponse
// @Failure 409 {object} errResponse
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Failure 501 {object} errResponse
// @Router /api/shorten/batch [post]
func (d *delivery) batchURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	userID, ok := r.Context().Value(ctxKeyUserID{}).(string)
	if !ok {
		d.logger.Error("invalid user id", ErrInvalidRequest,
			slog.String("userID", userID), slog.String("handler", "batchURL"))
		d.handelErrURL(w, r, ErrInvalidRequest)
		return
	}

	req := new([]batchURLRequest)
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		d.logger.Error("failed decode request url", err, slog.String("handler", "batchURL"))
		d.handelErrURL(w, r, ErrInvalidRequest)
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
			ErrInvalidRequest, slog.String("handler", "batchURL"))
		d.handelErrURL(w, r, ErrInvalidRequest)
		return
	}

	urls, err := d.shortener.BatchURL(r.Context(), correlationID, rawURL, userID)
	if err != nil {
		d.logger.Error("failed batch add urls", err, slog.String("handler", "batchURL"))
		d.handelErrURL(w, r, err)
		return
	}

	resp := make([]batchURLResponse, len(urls))

	for idx, url := range urls {
		resp[idx] = batchURLResponse{
			url.CorrelationID(),
			url.ShortURL(),
		}
	}

	data, err := json.Marshal(resp)
	if err != nil {
		d.logger.Error("failed marshal response url", err, slog.String("handler", "batchURL"))
		d.handelErrURL(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		d.logger.Error("write body", err, slog.String("handler", "batchURL"))
		d.handelErrURL(w, r, ErrInternalServer)
		return
	}
}

// deleteURL godoc
// @Summary remove multiple short URLs
// @Description remove multiple short URLs
// @ID deleteURL
// @Produce application/json
// @Param ids body []string true "short URL ids to delete"
// @Success 202
// @Failure 410 {object} errResponse
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Failure 501 {object} errResponse
// @Router /api/shorten/user/urls [delete]
func (d *delivery) deleteURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(ctxKeyUserID{}).(string)
	if !ok {
		d.logger.Error("invalid user id", ErrInvalidRequest,
			slog.String("userID", userID), slog.String("handler", "deleteBatchURL"))
		d.handelErrURL(w, r, ErrInvalidRequest)
		return
	}

	urls := new([]string)
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&urls); err != nil {
		d.logger.Error("failed decode request", err, slog.String("handler", "deleteURL"))
		d.handelErrURL(w, r, ErrInvalidRequest)
		return
	}

	err := d.shortener.DeleteURL(r.Context(), userID, *urls)
	if err != nil {
		d.logger.Error("failed delete urls", err, slog.String("handler", "deleteURL"))
		d.handelErrURL(w, r, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// ping godoc
// @Summary health check shortener storage
// @Description health check shortener storage
// @ID ping
// @Success 200
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /ping [get]
func (d *delivery) ping(w http.ResponseWriter, r *http.Request) {
	err := d.shortener.StorageCheck(r.Context())
	if err != nil {
		d.logger.Error("failed check storage", err, slog.String("handler", "ping"))
		d.handelErrURL(w, r, ErrStorageCheck)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// stats godoc
// @Summary shorturl statistics
// @Description shorturl statistics
// @ID stats
// @Success 200
// @Failure 400 {object} errResponse
// @Failure 403 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /internal/stats [get]
func (d *delivery) stats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(ctxKeyUserID{}).(string)
	if !ok {
		d.logger.Error("invalid user id", ErrInvalidRequest,
			slog.String("userID", userID), slog.String("handler", "stats"))
		d.handelErrURL(w, r, ErrInvalidRequest)
		return
	}

	stats, err := d.shortener.GetStats(r.Context())
	if err != nil {
		d.logger.Error("failed getting stats", err, slog.String("handler", "stats"))
		d.handelErrURL(w, r, ErrInternalServer)
		return
	}

	res := new(statsResponse)
	res.URLs = stats.URL().Count()
	res.Users = stats.User().Count()

	if err = json.NewEncoder(w).Encode(res); err != nil {
		d.logger.Error("failed encode response", err, slog.String("handler", "stats"))
		d.handelErrURL(w, r, ErrInternalServer)
		return
	}
}

func (d *delivery) handelErrURL(w http.ResponseWriter, r *http.Request, err error) {
	var httpStatus int

	switch {
	case errors.Is(err, shortener.ErrDecodeURL):
		httpStatus = http.StatusBadRequest
	case errors.Is(err, shortener.ErrParseURL):
		httpStatus = http.StatusBadRequest
	case errors.Is(err, shortener.ErrParseUUID):
		httpStatus = http.StatusBadRequest
	case errors.Is(err, ErrInvalidRequest):
		httpStatus = http.StatusBadRequest
	case errors.Is(err, entity.ErrNotFound):
		httpStatus = http.StatusNotFound
	case errors.Is(err, entity.ErrAlreadyExist):
		w.WriteHeader(http.StatusConflict)
		return
	case errors.Is(err, ErrStorageCheck):
		httpStatus = http.StatusInternalServerError
	case errors.Is(err, entity.ErrDeleted):
		httpStatus = http.StatusGone
	default:
		httpStatus = http.StatusNotImplemented
	}

	err = render.Render(w, r, errRender(httpStatus, err))
	if err != nil {
		d.logger.Error("go-chi render err", err)
	}
}
