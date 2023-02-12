package shortener

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sreway/shorturl/config"
	entity "github.com/sreway/shorturl/internal/domain/url"
	repoMock "github.com/sreway/shorturl/internal/usecases/adapters/storage/mocks"
)

func TestUseCase_CreateURL(t *testing.T) {
	baseURL := &url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
	}
	anyMock := gomock.Any()

	type args struct {
		rawURL  string
		counter uint64
	}

	tests := []struct {
		name    string
		args    args
		want    *entity.URL
		wantErr bool
	}{
		{
			name: "positive create url",
			args: args{
				rawURL:  "https://ya.ru",
				counter: 1000000000,
			},
			want: &entity.URL{
				ID: 1000000001,
				LongURL: &url.URL{
					Scheme: "https",
					Host:   "ya.ru",
				},
				ShortURL: &url.URL{
					Scheme: baseURL.Scheme,
					Host:   baseURL.Host,
					Path:   "15FTGh",
				},
			},
			wantErr: false,
		},
		{
			name: "negative create url (invalid raw url)",
			args: args{
				rawURL:  "invalid",
				counter: 1000000000,
			},
			want:    new(entity.URL),
			wantErr: true,
		},
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()
	repo := repoMock.NewMockURL(ctl)
	ctx := context.Background()

	for _, tt := range tests {
		cfg := config.ShortURL{
			BaseURL: baseURL,
			Counter: tt.args.counter,
		}
		uc := New(repo, &cfg)
		repo.EXPECT().Add(anyMock, anyMock, anyMock).Return(nil).AnyTimes()

		t.Run(tt.name, func(t *testing.T) {
			got, err := uc.CreateURL(ctx, tt.args.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.LongURL.Scheme, got.LongURL.Scheme)
			assert.Equal(t, tt.want.LongURL.Host, got.LongURL.Host)
			assert.Equal(t, tt.want.LongURL.Path, got.LongURL.Path)
			assert.Equal(t, tt.want.ShortURL.Scheme, got.ShortURL.Scheme)
			assert.Equal(t, tt.want.ShortURL.Host, got.ShortURL.Host)
			assert.Equal(t, tt.want.ShortURL.Path, got.ShortURL.Path)
		})
	}
}

func TestUseCase_GetURL(t *testing.T) {
	baseURL := &url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
	}
	anyMock := gomock.Any()
	errNotFoundMock := errors.New("not found")

	type (
		repoResp struct {
			url *url.URL
			err error
		}

		args struct {
			urlID    string
			counter  uint64
			repoResp repoResp
		}
	)

	tests := []struct {
		name    string
		args    args
		want    *entity.URL
		wantErr bool
	}{
		{
			name: "positive get url",
			args: args{
				urlID:   "15FTGh",
				counter: 1000000000,
				repoResp: repoResp{
					&url.URL{
						Scheme: "https",
						Host:   "ya.ru",
					},
					nil,
				},
			},
			want: &entity.URL{
				ID: 1000000001,
				LongURL: &url.URL{
					Scheme: "https",
					Host:   "ya.ru",
				},
				ShortURL: &url.URL{
					Scheme: baseURL.Scheme,
					Host:   baseURL.Host,
					Path:   "15FTGh",
				},
			},
			wantErr: false,
		},
		{
			name: "negative get url (invalid urlID)",
			args: args{
				urlID:   "invalid!",
				counter: 1000000000,
			},
			want:    new(entity.URL),
			wantErr: true,
		},

		{
			name: "negative get url (not found urlID)",
			args: args{
				urlID:   "15FTGk",
				counter: 1000000000,
				repoResp: repoResp{
					nil,
					errNotFoundMock,
				},
			},

			want:    new(entity.URL),
			wantErr: true,
		},
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()
	ctx := context.Background()

	for _, tt := range tests {
		repo := repoMock.NewMockURL(ctl)
		cfg := config.ShortURL{
			BaseURL: baseURL,
			Counter: tt.args.counter,
		}

		uc := New(repo, &cfg)

		repo.EXPECT().Get(anyMock, anyMock).Return(tt.args.repoResp.url, tt.args.repoResp.err).AnyTimes()

		t.Run(tt.name, func(t *testing.T) {
			got, err := uc.GetURL(ctx, tt.args.urlID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.LongURL.Scheme, got.LongURL.Scheme)
			assert.Equal(t, tt.want.LongURL.Host, got.LongURL.Host)
			assert.Equal(t, tt.want.LongURL.Path, got.LongURL.Path)
			assert.Equal(t, tt.want.ShortURL.Scheme, got.ShortURL.Scheme)
			assert.Equal(t, tt.want.ShortURL.Host, got.ShortURL.Host)
			assert.Equal(t, tt.want.ShortURL.Path, got.ShortURL.Path)
		})
	}
}
