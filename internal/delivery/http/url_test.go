package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/sreway/shorturl/internal/domain/url"
	urlMock "github.com/sreway/shorturl/internal/domain/url/mock"
	usecasesMock "github.com/sreway/shorturl/internal/usecases/mock"
	"github.com/sreway/shorturl/internal/usecases/shortener"
)

func Test_delivery_addURL(t *testing.T) {
	type want struct {
		code     int
		response string
		headers  map[string]string
	}
	type args struct {
		uri    string
		method string
		body   string
	}
	type fields struct {
		useCaseShortURL string
		useCaseErr      error
	}

	tests := []struct {
		name   string
		args   args
		fields fields
		want   want
	}{
		{
			name: "positive add url",
			args: args{
				uri:    "/",
				method: http.MethodPost,
				body:   `https://ya.ru`,
			},
			fields: fields{
				useCaseShortURL: "http://localhost:8080/2ZrI5IHFnvPscPYKlxFtRQ",
			},
			want: want{
				code:     http.StatusCreated,
				response: `http://localhost:8080/2ZrI5IHFnvPscPYKlxFtRQ`,
				headers: map[string]string{
					"Content-Type": "text/plain",
				},
			},
		},

		{
			name: "negative add url (invalid body)",
			args: args{
				uri:    "/",
				method: http.MethodPost,
				body:   `invalid`,
			},
			fields: fields{
				useCaseErr: shortener.ErrParseURL,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "{\"error\":\"URL parsing error\"}\n",
				headers: map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
			},
		},

		{
			name: "negative add url (empty body)",
			args: args{
				uri:    "/",
				method: http.MethodPost,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "{\"error\":\"invalid request\"}\n",
				headers: map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
			},
		},

		{
			name: "negative add url (exist url)",
			args: args{
				uri:    "/",
				method: http.MethodPost,
				body:   `https://ya.ru`,
			},
			fields: fields{
				useCaseShortURL: "http://localhost:8080/2ZrI5IHFnvPscPYKlxFtRQ",
				useCaseErr:      url.ErrAlreadyExist,
			},
			want: want{
				code:     http.StatusConflict,
				response: `http://localhost:8080/2ZrI5IHFnvPscPYKlxFtRQ`,
				headers: map[string]string{
					"Content-Type": "text/plain",
				},
			},
		},
	}

	anyMock := gomock.Any()
	userID := uuid.New().String()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	for _, tt := range tests {
		uc := usecasesMock.NewMockShortener(ctl)
		url := urlMock.NewMockURL(ctl)
		url.EXPECT().ShortURL().Return(tt.fields.useCaseShortURL).AnyTimes()
		uc.EXPECT().CreateURL(anyMock, anyMock, anyMock).Return(url, tt.fields.useCaseErr).AnyTimes()
		d := New(uc)
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.uri, strings.NewReader(tt.args.body))
			request = request.WithContext(context.WithValue(request.Context(), ctxKeyUserID{}, userID))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(d.addURL)
			h.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
			resBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			for k, v := range tt.want.headers {
				assert.Equal(t, resp.Header.Get(k), v)
			}
		})
	}
}

func Test_delivery_getURL(t *testing.T) {
	type want struct {
		code    int
		headers map[string]string
	}
	type args struct {
		uri    string
		method string
	}
	type fields struct {
		useCaseLongURL string
		useCaseErr     error
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   want
	}{
		{
			name: "positive get url",
			args: args{
				uri:    "/2ZrI5IHFnvPscPYKlxFtRQ",
				method: http.MethodGet,
			},
			fields: fields{
				useCaseLongURL: "https://ya.ru",
			},
			want: want{
				code: http.StatusTemporaryRedirect,
				headers: map[string]string{
					"Content-Type": "text/plain",
					"Location":     "https://ya.ru",
				},
			},
		},

		{
			name: "negative get url (invalid slug)",
			args: args{
				uri:    "/invalid!",
				method: http.MethodGet,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},

		{
			name: "negative get url (invalid id)",
			args: args{
				uri:    "/2ZrI5IHFnvPscPYKlxFtRQ2",
				method: http.MethodGet,
			},
			fields: fields{
				useCaseErr: shortener.ErrParseUUID,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},

		{
			name: "negative get url (not found)",
			args: args{
				uri:    "/2ZrI5IHFnvPscPYKlxFtRQ",
				method: http.MethodGet,
			},
			fields: fields{
				useCaseErr: url.ErrNotFound,
			},
			want: want{
				code: http.StatusNotFound,
			},
		},
	}

	anyMock := gomock.Any()
	userID := uuid.New().String()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	for _, tt := range tests {
		uc := usecasesMock.NewMockShortener(ctl)
		url := urlMock.NewMockURL(ctl)
		url.EXPECT().LongURL().Return(tt.fields.useCaseLongURL).AnyTimes()
		uc.EXPECT().GetURL(anyMock, anyMock).Return(url, tt.fields.useCaseErr).AnyTimes()
		d := New(uc)
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.uri, nil)
			request = request.WithContext(context.WithValue(request.Context(), ctxKeyUserID{}, userID))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(d.getURL)
			h.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
			for k, v := range tt.want.headers {
				assert.Equal(t, resp.Header.Get(k), v)
			}
		})
	}
}

func Test_delivery_shortURL(t *testing.T) {
	type want struct {
		code     int
		response string
		headers  map[string]string
	}
	type args struct {
		uri    string
		method string
		body   string
	}
	type fields struct {
		useCaseShortURL string
		useCaseErr      error
	}

	tests := []struct {
		name   string
		args   args
		fields fields
		want   want
	}{
		{
			name: "positive add url",
			args: args{
				uri:    "/api/shorten",
				method: http.MethodPost,
				body:   `{"url":"https://ya.ru"}`,
			},
			fields: fields{
				useCaseShortURL: "http://localhost:8080/2ZrI5IHFnvPscPYKlxFtRQ",
			},
			want: want{
				code:     http.StatusCreated,
				response: `{"result":"http://localhost:8080/2ZrI5IHFnvPscPYKlxFtRQ"}`,
				headers: map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
			},
		},

		{
			name: "negative add url (invalid url)",
			args: args{
				uri:    "/api/shorten",
				method: http.MethodPost,
				body:   `{"url":"invalid"}`,
			},
			fields: fields{
				useCaseErr: shortener.ErrDecodeURL,
			},
			want: want{
				code: http.StatusBadRequest,
				headers: map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
				response: "{\"error\":\"URL decoding error\"}\n",
			},
		},

		{
			name: "negative add url (invalid body)",
			args: args{
				uri:    "/api/shorten",
				method: http.MethodPost,
				body:   `invalid`,
			},
			want: want{
				code: http.StatusBadRequest,
				headers: map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
				response: "{\"error\":\"invalid request\"}\n",
			},
		},

		{
			name: "negative add url (exist url)",
			args: args{
				uri:    "/api/shorten",
				method: http.MethodPost,
				body:   `{"url":"https://ya.ru"}`,
			},
			fields: fields{
				useCaseShortURL: "http://localhost:8080/2ZrI5IHFnvPscPYKlxFtRQ",
				useCaseErr:      url.ErrAlreadyExist,
			},
			want: want{
				code:     http.StatusConflict,
				response: `{"result":"http://localhost:8080/2ZrI5IHFnvPscPYKlxFtRQ"}`,
				headers: map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
			},
		},
	}

	anyMock := gomock.Any()
	userID := uuid.New().String()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	for _, tt := range tests {
		uc := usecasesMock.NewMockShortener(ctl)
		url := urlMock.NewMockURL(ctl)
		url.EXPECT().ShortURL().Return(tt.fields.useCaseShortURL).AnyTimes()
		uc.EXPECT().CreateURL(anyMock, anyMock, anyMock).Return(url, tt.fields.useCaseErr).AnyTimes()
		d := New(uc)
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.uri, strings.NewReader(tt.args.body))
			request = request.WithContext(context.WithValue(request.Context(), ctxKeyUserID{}, userID))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(d.shortURL)
			h.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
			resBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			for k, v := range tt.want.headers {
				assert.Equal(t, resp.Header.Get(k), v)
			}
		})
	}
}

func Test_delivery_userURL(t *testing.T) {
	type want struct {
		code     int
		response string
		headers  map[string]string
	}
	type args struct {
		uri    string
		method string
	}
	type fields struct {
		useCaseURLs []struct {
			shortURL string
			longURL  string
		}
		useCaseErr error
	}

	tests := []struct {
		name   string
		args   args
		fields fields
		want   want
	}{
		{
			name: "positive get user urls",
			args: args{
				uri:    "/api/user/urls",
				method: http.MethodGet,
			},
			fields: fields{
				useCaseURLs: []struct {
					shortURL string
					longURL  string
				}{
					{
						shortURL: "http://localhost:8080/2ShKzidROaM6mhK2RP7chv",
						longURL:  "https://ya.ru",
					},
				},
			},
			want: want{
				code:     http.StatusOK,
				response: `[{"short_url":"http://localhost:8080/2ShKzidROaM6mhK2RP7chv","original_url":"https://ya.ru"}]`,
				headers: map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
			},
		},

		{
			name: "positive get user urls (no urls)",
			args: args{
				uri:    "/api/user/urls",
				method: http.MethodGet,
			},
			want: want{
				code: http.StatusNoContent,
				headers: map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
			},
		},

		{
			name: "negative get user urls (invalid userID)",
			args: args{
				uri:    "/api/user/urls",
				method: http.MethodGet,
			},
			fields: fields{
				useCaseErr: shortener.ErrParseUUID,
			},
			want: want{
				code: http.StatusBadRequest,
				headers: map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
				response: "{\"error\":\"UUID parsing error\"}\n",
			},
		},
	}

	anyMock := gomock.Any()
	userID := uuid.New().String()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	for _, tt := range tests {
		uc := usecasesMock.NewMockShortener(ctl)
		urls := []url.URL{}
		for _, item := range tt.fields.useCaseURLs {
			entity := urlMock.NewMockURL(ctl)
			entity.EXPECT().ShortURL().Return(item.shortURL).AnyTimes()
			entity.EXPECT().LongURL().Return(item.longURL).AnyTimes()
			urls = append(urls, entity)
		}
		uc.EXPECT().GetUserURLs(anyMock, anyMock).Return(urls, tt.fields.useCaseErr).AnyTimes()
		d := New(uc)
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.uri, nil)
			request = request.WithContext(context.WithValue(request.Context(), ctxKeyUserID{}, userID))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(d.userURL)
			h.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
			resBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			for k, v := range tt.want.headers {
				assert.Equal(t, resp.Header.Get(k), v)
			}
		})
	}
}

func Test_delivery_batchURL(t *testing.T) {
	type want struct {
		code     int
		response string
		headers  map[string]string
	}
	type args struct {
		uri    string
		method string
		body   string
	}
	type fields struct {
		useCaseURLs []struct {
			shortURL      string
			correlationID string
		}
		useCaseErr error
	}

	tests := []struct {
		name   string
		args   args
		fields fields
		want   want
	}{
		{
			name: "positive add batch url",
			args: args{
				uri:    "/api/shorten/batch",
				method: http.MethodPost,
				body:   `[{"correlation_id":"1","original_url":"https://ya.ru"}]`,
			},
			fields: fields{
				useCaseURLs: []struct {
					shortURL      string
					correlationID string
				}{
					{
						shortURL:      "http://localhost:8080/2ShKzidROaM6mhK2RP7chv",
						correlationID: "1",
					},
				},
			},
			want: want{
				code:     http.StatusCreated,
				response: `[{"correlation_id":"1","short_url":"http://localhost:8080/2ShKzidROaM6mhK2RP7chv"}]`,
				headers: map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
			},
		},

		{
			name: "negative add batch url (incorrect body)",
			args: args{
				uri:    "/api/shorten/batch",
				method: http.MethodPost,
				body:   `[{"correlation_id":"1","original_url":"https://ya.ru"},{"correlation_id":"2"}]`,
			},
			want: want{
				code: http.StatusBadRequest,
				headers: map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
				response: "{\"error\":\"invalid request\"}\n",
			},
		},

		{
			name: "negative add batch url (invalid body)",
			args: args{
				uri:    "/api/shorten/batch",
				method: http.MethodPost,
				body:   `invalid`,
			},
			want: want{
				code: http.StatusBadRequest,
				headers: map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
				response: "{\"error\":\"invalid request\"}\n",
			},
		},

		{
			name: "negative add batch url (exist url)",
			args: args{
				uri:    "/api/shorten/batch",
				method: http.MethodPost,
				body:   `[{"correlation_id":"1","original_url":"https://ya.ru"}]`,
			},
			fields: fields{
				useCaseURLs: []struct {
					shortURL      string
					correlationID string
				}{
					{
						shortURL:      "http://localhost:8080/7jUSM5AZgNW4DwdHts5b0F",
						correlationID: "1",
					},
				},
				useCaseErr: url.ErrAlreadyExist,
			},
			want: want{
				code: http.StatusConflict,
				headers: map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
			},
		},
	}

	anyMock := gomock.Any()
	userID := uuid.New().String()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	for _, tt := range tests {
		uc := usecasesMock.NewMockShortener(ctl)
		urls := []url.URL{}
		for _, item := range tt.fields.useCaseURLs {
			entity := urlMock.NewMockURL(ctl)
			entity.EXPECT().ShortURL().Return(item.shortURL).AnyTimes()
			entity.EXPECT().CorrelationID().Return(item.correlationID).AnyTimes()
			urls = append(urls, entity)
		}
		uc.EXPECT().BatchURL(anyMock, anyMock, anyMock, anyMock).Return(urls, tt.fields.useCaseErr).AnyTimes()
		d := New(uc)
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.uri, strings.NewReader(tt.args.body))
			request = request.WithContext(context.WithValue(request.Context(), ctxKeyUserID{}, userID))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(d.batchURL)
			h.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
			resBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			for k, v := range tt.want.headers {
				assert.Equal(t, resp.Header.Get(k), v)
			}
		})
	}
}

func Test_delivery_ping(t *testing.T) {
	type want struct {
		code int
	}
	type args struct {
		uri    string
		method string
	}
	type fields struct {
		useCaseErr error
	}

	tests := []struct {
		name   string
		args   args
		fields fields
		want   want
	}{
		{
			name: "positive ping",
			args: args{
				uri:    "/ping",
				method: http.MethodGet,
			},
			want: want{
				code: http.StatusOK,
			},
		},

		{
			name: "negative url",
			args: args{
				uri:    "/ping",
				method: http.MethodGet,
			},
			fields: fields{
				useCaseErr: ErrStorageCheck,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}

	anyMock := gomock.Any()
	userID := uuid.New().String()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	for _, tt := range tests {
		uc := usecasesMock.NewMockShortener(ctl)
		uc.EXPECT().StorageCheck(anyMock).Return(tt.fields.useCaseErr).AnyTimes()
		d := New(uc)
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.uri, nil)
			request = request.WithContext(context.WithValue(request.Context(), ctxKeyUserID{}, userID))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(d.ping)
			h.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
		})
	}
}

func Test_delivery_deleteURL(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	type args struct {
		body string
	}
	type fields struct {
		useCaseErr error
	}

	uri := "/api/user/urls"
	method := http.MethodDelete

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "positive delete url",
			args: args{
				body: `["2ZrI5IHFnvPscPYKlxFtRQ"]`,
			},
			want: want{
				code: http.StatusAccepted,
			},
		},
		{
			name: "negative delete url (invalid body)",
			args: args{
				body: `invalid`,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "{\"error\":\"invalid request\"}\n",
			},
		},
		{
			name: "negative delete url (invalid url id, not uuid)",
			args: args{
				body: `["2ZrI5IHFnvPscPYKlxFtRQs"]`,
			},
			fields: fields{useCaseErr: shortener.ErrParseUUID},
			want: want{
				code:     http.StatusBadRequest,
				response: "{\"error\":\"UUID parsing error\"}\n",
			},
		},
		{
			name: "negative delete url (already deleted)",
			args: args{
				body: `["2ZrI5IHFnvPscPYKlxFtRQ"]`,
			},
			fields: fields{useCaseErr: url.ErrDeleted},
			want: want{
				code:     http.StatusGone,
				response: "{\"error\":\"URL deleted\"}\n",
			},
		},
	}

	anyMock := gomock.Any()
	userID := uuid.New().String()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	for _, tt := range tests {
		uc := usecasesMock.NewMockShortener(ctl)
		uc.EXPECT().DeleteURL(anyMock, anyMock, anyMock).Return(tt.fields.useCaseErr).AnyTimes()
		d := New(uc)
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(method, uri, strings.NewReader(tt.args.body))
			request = request.WithContext(context.WithValue(request.Context(), ctxKeyUserID{}, userID))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(d.deleteURL)
			h.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
			resBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
		})
	}
}
