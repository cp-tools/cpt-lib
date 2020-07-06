package codeforces

import (
	"net/http"
	"net/http/cookiejar"
	"os"
	"reflect"
	"testing"
)

func init() {
	// login to account for access to all other tests
	usr := os.Getenv("CODEFORCES_USERNAME")
	passwd := os.Getenv("CODEFORCES_PASSWORD")

	jar, _ := cookiejar.New(nil)
	SessCln = &http.Client{Jar: jar}
	Login(usr, passwd)
}

func Test_loginPage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Login Page",
			want: "https://codeforces.com/enter",
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

func TestLogin(t *testing.T) {
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
			name: "Login to cp-tools account",
			args: args{
				os.Getenv("CODEFORCES_USERNAME"),
				os.Getenv("CODEFORCES_PASSWORD"),
			},
			want:    "cp-tools",
			wantErr: false,
		},
		{
			name:    "Invalid login",
			args:    args{"infixint943", "ThIsNoTmYPASsWd"},
			want:    "",
			wantErr: true,
		},
		// TODO: Add test cases.
	}

	tmpSess := *SessCln
	for _, tt := range tests {
		// reset cookie configurations
		jar, _ := cookiejar.New(nil)
		SessCln = &http.Client{Jar: jar}

		t.Run(tt.name, func(t *testing.T) {
			got, err := Login(tt.args.usr, tt.args.passwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
		})
	}
	SessCln = &tmpSess
}