package shortener

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_encodeUUID(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive encode uuid",
			args: args{
				uuid: "624708fa-d258-4b99-b09a-49d95f294626",
			},
			want: "2ZrI5IHFnvPscPYKlxFtRQ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := uuid.Parse(tt.args.uuid)
			assert.NoError(t, err)
			if got := encodeUUID(id); got != tt.want {
				t.Errorf("encodeUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeUUID(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positive decode uuid",
			args: args{
				s: "2ZrI5IHFnvPscPYKlxFtRQ",
			},
			want:    "624708fa-d258-4b99-b09a-49d95f294626",
			wantErr: false,
		},

		{
			name: "negative decode uuid",
			args: args{
				s: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeUUID(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				return
			}

			id, err := uuid.FromBytes(got)
			assert.NoError(t, err)

			if id.String() != tt.want {
				t.Errorf("decodeUUID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
