package codeforces

import (
	"reflect"
	"testing"
	"time"
)

func Test_parseTime(t *testing.T) {
	type args struct {
		link string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "Test #1",
			args: args{"Oct/20/2020 20:05 UTC+3.0"},
			want: time.Date(2020, time.October, 20, 17, 5, 0, 0, time.UTC),
		},
		{
			name: "Test #2",
			args: args{"Mar/13/2022 14:05"},
			want: time.Date(2022, time.March, 13, 14, 5, 0, 0, time.UTC),
		},
		{
			name: "Test #3", // Russian locale.
			args: args{"23.10.2020 15:35 UTC-2.0"},
			want: time.Date(2020, time.October, 23, 17, 35, 0, 0, time.UTC),
		},
		{
			name: "Test #4", // Russian locale.
			args: args{"23.10.2020 15:35"},
			want: time.Date(2020, time.October, 23, 15, 35, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTime(tt.args.link); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
