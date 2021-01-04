package atcoder

import (
	"testing"
	"time"
)

func TestArgs_GetCountdown(t *testing.T) {
	tests := []struct {
		name    string
		arg     Args
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "Test #1",
			arg:     Args{"abc180", ""},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Test #2",
			arg:     Args{"InVaLiD123", ""},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.GetCountdown()
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.GetCountdown() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Args.GetCountdown() = %v, want %v", got, tt.want)
			}
		})
	}
}
