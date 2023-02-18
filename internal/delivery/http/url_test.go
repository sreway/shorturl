package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	entity "github.com/sreway/shorturl/internal/domain/url"
	ucMock "github.com/sreway/shorturl/internal/usecases/mocks"
	"github.com/sreway/shorturl/internal/usecases/shortener"
)

func TestDelivery_addURL(t *testing.T) {
	type (
		want struct {
			code     int
			response string
			headers  map[string]string
		}

		ucResp struct {
			url *entity.URL
			err error
		}

		args struct {
			body   string
			method string
			ucResp ucResp
		}
	)

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive add url",
			args: args{
				body:   `https://ya.ru`,
				method: http.MethodPost,
				ucResp: ucResp{
					&entity.URL{
						ShortURL: &url.URL{
							Scheme: "http",
							Host:   "localhost:8080",
							Path:   "15FTGh",
						},
					},
					nil,
				},
			},
			want: want{
				code:     http.StatusCreated,
				response: `http://localhost:8080/15FTGh`,
				headers: map[string]string{
					"Content-Type": "text/plain",
				},
			},
		},
		{
			name: "negative add url (invalid body)",
			args: args{
				body:   `invalid`,
				method: http.MethodPost,
				ucResp: ucResp{
					nil,
					shortener.ErrParseURL,
				},
			},
			want: want{
				code: http.StatusBadRequest,
				headers: map[string]string{
					"Content-Type": "text/plain",
				},
			},
		},

		{
			name: "negative add url (empty body)",
			args: args{
				method: http.MethodPost,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	anyMock := gomock.Any()

	for _, tt := range tests {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		uc := ucMock.NewMockShortener(ctl)
		d := New(uc)
		uc.EXPECT().CreateURL(anyMock, anyMock).Return(tt.args.ucResp.url, tt.args.ucResp.err).AnyTimes()
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, "/", strings.NewReader(tt.args.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(d.addURL)
			h.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
			resBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, string(resBody), tt.want.response)

			for k, v := range tt.want.headers {
				assert.Equal(t, resp.Header.Get(k), v)
			}
		})
	}
}

func TestDelivery_getURL(t *testing.T) {
	type (
		want struct {
			code    int
			headers map[string]string
		}

		ucResp struct {
			url *entity.URL
			err error
		}

		args struct {
			method string
			path   string
			ucResp ucResp
		}
	)

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive get url",
			args: args{
				path:   `/15FTGh`,
				method: http.MethodGet,
				ucResp: ucResp{
					&entity.URL{
						LongURL: &url.URL{
							Scheme: "https",
							Host:   "ya.ru",
						},
					},
					nil,
				},
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
				path:   `/invalid!`,
				method: http.MethodGet,
				ucResp: ucResp{
					nil,
					ErrInvalidSlug,
				},
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},

		{
			name: "negative get url (invalid id)",
			args: args{
				path:   `/-15FTGh!`,
				method: http.MethodGet,
				ucResp: ucResp{
					nil,
					shortener.ErrDecodeURL,
				},
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	anyMock := gomock.Any()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	for _, tt := range tests {
		uc := ucMock.NewMockShortener(ctl)
		d := New(uc)
		uc.EXPECT().GetURL(anyMock, anyMock).Return(tt.args.ucResp.url, tt.args.ucResp.err).AnyTimes()
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.path, nil)
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
