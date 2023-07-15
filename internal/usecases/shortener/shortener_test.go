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
				repoErr: url.NewURLErr(uuid.New(), uuid.New(), url.ErrAlreadyExist),
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
		uc := New(repo, cfg.GetShortURL())

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
		urlID   string
		deleted bool
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
				urlID:   "5nPymsbLZfXlsUDlZ4MIhY",
				deleted: false,
			},
			wantErr: assert.NoError,
		},

		{
			name: "negative get url (not found)",
			args: args{
				urlID:   "5nPymsbLZfXlsUDlZ4MIhY",
				deleted: false,
			},
			fields: fields{
				repoErr: url.ErrNotFound,
			},
			wantErr: assert.Error,
		},

		{
			name: "negative get url (invalid url id)",
			args: args{
				urlID:   "invalid",
				deleted: false,
			},
			wantErr: assert.Error,
		},

		{
			name: "negative get url (invalid url id, not uuid)",
			args: args{
				urlID:   "2ZrI5IHFnvPscPYKlxFtRQs",
				deleted: false,
			},
			wantErr: assert.Error,
		},

		{
			name: "negative get url (deleted)",
			args: args{
				urlID:   "5nPymsbLZfXlsUDlZ4MIhY",
				deleted: true,
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
		uc := New(repo, cfg.GetShortURL())
		mockURL := urlMock.NewMockURL(ctl)
		mockURL.EXPECT().Deleted().Return(tt.args.deleted).AnyTimes()
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
			name: "negative get urls (not found)",
			args: args{
				userID: "invalid",
			},
			fields: fields{
				repoErr: url.ErrNotFound,
			},
			wantErr: assert.Error,
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
		uc := New(repo, cfg.GetShortURL())
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
		uc := New(repo, cfg.GetShortURL())
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
		uc := New(repo, cfg.GetShortURL())
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

func Test_useCase_DeleteURL(t *testing.T) {
	type args struct {
		userID string
		urlID  []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive delete url",
			args: args{
				userID: "035f67d8-626b-48f2-b436-8509954fc452",
				urlID:  []string{"5nPymsbLZfXlsUDlZ4MIhY"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "negative delete url (invalid user uuid)",
			args: args{
				userID: "invalid",
			},
			wantErr: assert.Error,
		},
		{
			name: "negative delete url (invalid url id)",
			args: args{
				userID: "035f67d8-626b-48f2-b436-8509954fc452",
				urlID:  []string{"5nPymsbLZfXlsUDlZ4MIhY", "invalid"},
			},
			wantErr: assert.Error,
		},
		{
			name: "negative delete url (invalid url id, not uuid)",
			args: args{
				userID: "035f67d8-626b-48f2-b436-8509954fc452",
				urlID:  []string{"2ZrI5IHFnvPscPYKlxFtRQs"},
			},
			wantErr: assert.Error,
		},
	}
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	for _, tt := range tests {
		cfg, err := config.NewConfig()
		assert.NoError(t, err)
		repo := repoMock.NewMockURL(ctl)
		uc := New(repo, cfg.GetShortURL())
		t.Run(tt.name, func(t *testing.T) {
			err = uc.DeleteURL(ctx, tt.args.userID, tt.args.urlID)
			if !tt.wantErr(t, err, fmt.Sprintf("DeleteURL(%v)", tt.args.urlID)) {
				return
			}
			if err != nil {
				return
			}
		})
	}
}

func Test_useCase_GetStats(t *testing.T) {
	type fields struct {
		urlCount  int
		userCount int
		repoErr   error
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive get stats",
			fields: fields{
				urlCount:  5,
				userCount: 4,
			},
			wantErr: assert.NoError,
		},
		{
			name: "negative get stats",
			fields: fields{
				repoErr: errors.New("some error"),
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
		repo.EXPECT().GetUserCount(anyMock).Return(tt.fields.userCount, tt.fields.repoErr).AnyTimes()
		repo.EXPECT().GetURLCount(anyMock).Return(tt.fields.urlCount, tt.fields.repoErr).AnyTimes()
		uc := New(repo, cfg.GetShortURL())
		t.Run(tt.name, func(t *testing.T) {
			got, err := uc.GetStats(ctx)
			if tt.wantErr(t, err, "GetStat") {
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.fields.userCount, got.User().Count())
			assert.Equal(t, tt.fields.urlCount, got.URL().Count())
		})
	}
}
