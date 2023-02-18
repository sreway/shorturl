package shortener

import (
	"testing"
)

func Test_uintEncode(t *testing.T) {
	type args struct {
		number uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive encode",
			args: args{
				number: 1000000001,
			},
			want: "15FTGh",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uintEncode(tt.args.number); got != tt.want {
				t.Errorf("uintEncode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uintDecode(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "positive decode",
			args: args{
				"15FTGh",
			},
			want:    1000000001,
			wantErr: false,
		},

		{
			name: "negative decode (incorrect string)",
			args: args{
				"incorrect!",
			},
			wantErr: true,
		},

		{
			name: "negative decode (not uint64)",
			args: args{
				"-15FTGg",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uintDecode(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("uintDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("uintDecode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
