package codeforces

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	// setup headless browser to use
	_, mode := os.LookupEnv("LOCAL_MODE")

	Start(!mode, "../user-data-dir", "google-chrome",
		[]string{"blink-settings", "imagesEnabled=false"})

	if !mode {
		// setup login access to use
		usr := os.Getenv("CODEFORCES_USERNAME")
		passwd := os.Getenv("CODEFORCES_PASSWORD")
		_, err := login(usr, passwd)
		if err != nil {
			fmt.Println("Login failed:", err)
			Browser.Close()
			os.Exit(1)
		}
	}

	exitCode := m.Run()

	// logout current user
	if !mode {
		logout()
	}

	Browser.Close()
	os.Exit(exitCode)
}

func TestParse(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    Args
		wantErr bool
	}{
		{
			name:    "(URL)Test #1",
			args:    args{"https://codeforces.com/contest/1355"},
			want:    Args{"1355", "", "contest", ""},
			wantErr: false,
		},
		{
			name:    "(URL)Test #2",
			args:    args{"https://codeforces.com/contest/739/problem/E"},
			want:    Args{"739", "e", "contest", ""},
			wantErr: false,
		},
		{
			name:    "(URL)Test #3",
			args:    args{"https://codeforces.com/gym/102595"},
			want:    Args{"102595", "", "gym", ""},
			wantErr: false,
		},
		{
			name:    "(URL)Test #4",
			args:    args{"https://codeforces.com/gym/102302/problem/i"},
			want:    Args{"102302", "i", "gym", ""},
			wantErr: false,
		},
		{
			name:    "(URL)Test #5",
			args:    args{"https://codeforces.com/group/MEqF8b6wBT/contest/277493"},
			want:    Args{"277493", "", "group", "MEqF8b6wBT"},
			wantErr: false,
		},
		{
			name:    "(URL)Test #6",
			args:    args{"https://codeforces.com/group/MEqF8b6wBT/contest/277493/problem/g"},
			want:    Args{"277493", "g", "group", "MEqF8b6wBT"},
			wantErr: false,
		},
		{
			name:    "(UFO)Test #1",
			args:    args{"1355"},
			want:    Args{"1355", "", "contest", ""},
			wantErr: false,
		},
		{
			name:    "(UFO)Test #2",
			args:    args{"739 E"},
			want:    Args{"739", "e", "contest", ""},
			wantErr: false,
		},
		{
			name:    "(UFO)Test #3",
			args:    args{"102595"},
			want:    Args{"102595", "", "gym", ""},
			wantErr: false,
		},
		{
			name:    "(UFO)Test #4",
			args:    args{"102302i"},
			want:    Args{"102302", "i", "gym", ""},
			wantErr: false,
		},
		{
			name:    "(UFO)Test #5",
			args:    args{"MEqF8b6wBT 277493"},
			want:    Args{"277493", "", "group", "MEqF8b6wBT"},
			wantErr: false,
		},
		{
			name:    "(UFO)Test #6",
			args:    args{"MEqF8b6wBT 277493 g"},
			want:    Args{"277493", "g", "group", "MEqF8b6wBT"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
