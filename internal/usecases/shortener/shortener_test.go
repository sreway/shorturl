package shortener

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/sreway/shorturl/internal/config"
	"github.com/sreway/shorturl/internal/domain/url"
	urlMock "github.com/sreway/shorturl/internal/domain/url/mock"
	repoMock "github.com/sreway/shorturl/internal/usecases/adapters/storage/mock"
)

func Test_useCase_CreateURL(t *testing.T) {
	type args struct {
		rawURL string
		userID string
	}
	type fields struct {
		repoErr error
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive create url",
			args: args{
				rawURL: "https://ya.ru",
				userID: "624708fa-d258-4b99-b09a-49d95f294626",
			},
			wantErr: assert.NoError,
		},

		{
			name: "negative create url (invalid url)",
			args: args{
				rawURL: "invalid",
				userID: "624708fa-d258-4b99-b09a-49d95f294626",
			},
			wantErr: assert.Error,
		},

		{
			name: "negative create url (invalid user id)",
			args: args{
				rawURL: "https://ya.ru",
				userID: "invalid",
			},
			wantErr: assert.Error,
		},

		{
			name: "negative create url (exist url)",
			args: args{
				rawURL: "https://ya.ru",
				userID: "624708fa-d258-4b99-b09a-49d95f294626",
			},
			fields: fields{
				repoErr: url.ErrAlreadyExist,
			},
			wantErr: assert.Error,
		},
	}
	anyMock := gomock.Any()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	for _, tt := range tests {
		cfg, err := config.NewConfig()
		assert.NoError(t, err)
		repo := repoMock.NewMockURL(ctl)
		uc := New(repo, cfg.ShortURL())
		repo.EXPECT().Add(anyMock, anyMock).Return(tt.fields.repoErr).AnyTimes()
		t.Run(tt.name, func(t *testing.T) {
			got, err := uc.CreateURL(ctx, tt.args.rawURL, tt.args.userID)
			if !tt.wantErr(t, err, fmt.Sprintf("CreateURL(%v, %v)", tt.args.rawURL, tt.args.userID)) {
				return
			}
			if err != nil {
				return
			}

			assert.Equal(t, got.UserID().String(), tt.args.userID)
			assert.Equal(t, got.LongURL(), tt.args.rawURL)
			assert.NotEmpty(t, got.ShortURL())
			assert.NotEmpty(t, got.ID().String())
		})
	}
}

func Test_useCase_GetURL(t *testing.T) {
	type args struct {
		urlID string
	}
	type fields struct {
		repoErr error
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive get url",
			args: args{
				urlID: "5nPymsbLZfXlsUDlZ4MIhY",
			},
			wantErr: assert.NoError,
		},

		{
			name: "negative get url (not found)",
			args: args{
				urlID: "5nPymsbLZfXlsUDlZ4MIhY",
			},
			fields: fields{
				repoErr: url.ErrNotFound,
			},
			wantErr: assert.Error,
		},

		{
			name: "negative get url (invalid url id)",
			args: args{
				urlID: "invalid",
			},
			wantErr: assert.Error,
		},
	}
	anyMock := gomock.Any()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	for _, tt := range tests {
		cfg, err := config.NewConfig()
		assert.NoError(t, err)
		repo := repoMock.NewMockURL(ctl)
		uc := New(repo, cfg.ShortURL())
		mockURL := urlMock.NewMockURL(ctl)
		mockURL.EXPECT().SetShortURL(anyMock).AnyTimes()
		repo.EXPECT().Get(anyMock, anyMock).Return(mockURL, tt.fields.repoErr).AnyTimes()
		t.Run(tt.name, func(t *testing.T) {
			_, err = uc.GetURL(ctx, tt.args.urlID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetURL(%v)", tt.args.urlID)) {
				return
			}
			if err != nil {
				return
			}
		})
	}
}

func Test_useCase_GetUserURLs(t *testing.T) {
	type args struct {
		userID string
	}

	type fields struct {
		repoErr error
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive get user urls",
			args: args{
				userID: "035f67d8-626b-48f2-b436-8509954fc452",
			},
			wantErr: assert.NoError,
		},

		{
			name: "negative get urls (invalid user id)",
			args: args{
				userID: "invalid",
			},
			wantErr: assert.Error,
		},
	}
	anyMock := gomock.Any()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	for _, tt := range tests {
		cfg, err := config.NewConfig()
		assert.NoError(t, err)
		repo := repoMock.NewMockURL(ctl)
		uc := New(repo, cfg.ShortURL())
		urls := []url.URL{}
		mockURL := urlMock.NewMockURL(ctl)
		mockURL.EXPECT().SetShortURL(anyMock).AnyTimes()
		mockURL.EXPECT().ID().Return(uuid.New()).AnyTimes()
		urls = append(urls, mockURL)

		repo.EXPECT().GetByUserID(anyMock, anyMock).Return(urls, tt.fields.repoErr).AnyTimes()
		t.Run(tt.name, func(t *testing.T) {
			_, err = uc.GetUserURLs(ctx, tt.args.userID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetUserURLs(%v)", tt.args.userID)) {
				return
			}
		})
	}
}

func Test_useCase_StorageCheck(t *testing.T) {
	type fields struct {
		repoErr error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "positive check storage",
			wantErr: assert.NoError,
		},
		{
			name: "negative check storage",
			fields: fields{
				repoErr: errors.New("any error"),
			},
			wantErr: assert.Error,
		},
	}
	anyMock := gomock.Any()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	for _, tt := range tests {
		cfg, err := config.NewConfig()
		assert.NoError(t, err)
		repo := repoMock.NewMockURL(ctl)
		uc := New(repo, cfg.ShortURL())
		repo.EXPECT().Ping(anyMock).Return(tt.fields.repoErr).AnyTimes()
		t.Run(tt.name, func(t *testing.T) {
			err = uc.StorageCheck(ctx)
			if !tt.wantErr(t, err, "StorageCheck()") {
				return
			}
		})
	}
}

func Test_useCase_BatchURL(t *testing.T) {
	type args struct {
		userID        string
		rawURL        []string
		correlationID []string
	}
	type fields struct {
		repoErr error
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive batch url",
			args: args{
				userID: "035f67d8-626b-48f2-b436-8509954fc452",
				rawURL: []string{
					"https://ya.ru",
				},
				correlationID: []string{
					"1",
				},
			},
			wantErr: assert.NoError,
		},

		{
			name: "negative batch url (invalid user id)",
			args: args{
				userID: "invalid",
				rawURL: []string{
					"https://ya.ru",
				},
				correlationID: []string{
					"1",
				},
			},
			wantErr: assert.Error,
		},

		{
			name: "negative batch url (invalid raw url)",
			args: args{
				userID: "035f67d8-626b-48f2-b436-8509954fc452",
				rawURL: []string{
					"invalid",
				},
				correlationID: []string{
					"1",
				},
			},
			wantErr: assert.Error,
		},

		{
			name: "negative batch url (exist url)",
			args: args{
				userID: "035f67d8-626b-48f2-b436-8509954fc452",
				rawURL: []string{
					"https://ya.ru",
				},
				correlationID: []string{
					"1",
				},
			},
			fields: fields{
				repoErr: url.ErrAlreadyExist,
			},
			wantErr: assert.Error,
		},
	}
	anyMock := gomock.Any()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	for _, tt := range tests {
		cfg, err := config.NewConfig()
		assert.NoError(t, err)
		repo := repoMock.NewMockURL(ctl)
		uc := New(repo, cfg.ShortURL())
		repo.EXPECT().Batch(anyMock, anyMock).Return(tt.fields.repoErr).AnyTimes()
		t.Run(tt.name, func(t *testing.T) {
			got, err := uc.BatchURL(ctx, tt.args.correlationID, tt.args.rawURL, tt.args.userID)
			if !tt.wantErr(t, err, fmt.Sprintf("BatchURL(%v, %v)", tt.args.rawURL, tt.args.userID)) {
				return
			}
			if err != nil && !errors.Is(err, url.ErrAlreadyExist) {
				return
			}

			for _, item := range got {
				assert.NotEmpty(t, item.ShortURL())
				assert.NotEmpty(t, item.LongURL())
				assert.NotEmpty(t, item.ID().String())
				assert.NotEmpty(t, item.UserID().String())
				assert.NotEmpty(t, item.CorrelationID())
			}
		})
	}
}
