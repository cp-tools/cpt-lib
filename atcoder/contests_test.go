package atcoder

import (
	"testing"
	"time"
)

func TestArgs_GetCountdown(t *testing.T) {
	type fields struct {
		isVC bool
	}
	tests := []struct {
		name    string
		arg     Args
		fields  fields
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "Test #1",
			fields:  fields{true},
			arg:     Args{"abc180", ""},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.GetCountdown(tt.fields.isVC)
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
