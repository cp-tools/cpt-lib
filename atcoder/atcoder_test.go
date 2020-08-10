package atcoder

import (
	"os"
	"reflect"
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func init() {
	url := launcher.New().UserDataDir("user-data-dir").
		Set("blink-settings", "imagesEnabled=false").Launch()
	Browser = rod.New().ControlURL(url).Connect()

	// login to account for access to all other tests
	usr := os.Getenv("ATCODER_USERNAME")
	passwd := os.Getenv("ATCODER_PASSWORD")

	login(usr, passwd)
}

func Test_loginPage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Test #1",
			want: "https://atcoder.jp/login",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loginPage(); got != tt.want {
				t.Errorf("loginPage() = %v, want %v", got, tt.want)
			}
		})
	}
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
			name:    "Test #1",
			args:    args{"https://atcoder.jp/contests/m-solutions2020"},
			want:    Args{"m-solutions2020", ""},
			wantErr: false,
		},
		{
			name:    "Test #2",
			args:    args{"atcoder.jp/contests/m-solutions2020/tasks/m_solutions2020_a"},
			want:    Args{"m-solutions2020", "m_solutions2020_a"},
			wantErr: false,
		},
		{
			name:    "Test #3",
			args:    args{"atcoder.jp/contests/m_solutions2020"},
			want:    Args{},
			wantErr: true,
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

func Test_login(t *testing.T) {
	type args struct {
		usr    string
		passwd string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Login to cptools",
			args: args{
				os.Getenv("ATCODER_USERNAME"),
				os.Getenv("ATCODER_PASSWORD"),
			},
			want:    "cptools",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := login(tt.args.usr, tt.args.passwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("login() = %v, want %v", got, tt.want)
			}
		})
	}
}
