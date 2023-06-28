package http

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sreway/shorturl/internal/config"
	mockStorage "github.com/sreway/shorturl/internal/usecases/adapters/storage/mock"
	"github.com/sreway/shorturl/internal/usecases/shortener"
)

func Test_delivery_Run1(t *testing.T) {
	type args struct {
		env map[string]string
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "negative run (invalid server address)",
			args: args{
				env: map[string]string{
					"SERVER_ADDRESS": "invalid",
				},
				ctx: context.Background(),
			},
			wantErr: assert.Error,
		},
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.args.env {
				err := os.Setenv(k, v)
				assert.NoError(t, err)
			}
			cfg, err := config.NewConfig()
			assert.NoError(t, err)

			repo := mockStorage.NewMockURL(ctl)

			uc := shortener.New(repo, cfg.GetShortURL())
			d := New(uc)
			tt.wantErr(t, d.Run(tt.args.ctx, cfg.GetHTTP()), fmt.Sprintf("Run(%v, %v)", tt.args.ctx, cfg))

			for k := range tt.args.env {
				err = os.Unsetenv(k)
				assert.NoError(t, err)
			}
		})
	}
}
