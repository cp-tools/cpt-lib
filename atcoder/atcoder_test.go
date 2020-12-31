package atcoder

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
)

func getLoginCredentials() (string, string) {
	// setup login access to use
	usr := os.Getenv("ATCODER_USERNAME")
	passwd := os.Getenv("ATCODER_PASSWORD")
	return usr, passwd
}

func TestMain(m *testing.M) {
	// Load local .env file.
	godotenv.Load()

	_, browserHeadless := os.LookupEnv("BROWSER_HEADLESS")
	browserBin := os.Getenv("BROWSER_BINARY")
	if err := Start(browserHeadless, "", browserBin); err != nil {
		fmt.Println("Failed to start browser:", err)
		os.Exit(1)
	}

	if _, err := login(getLoginCredentials()); err != nil {
		fmt.Println("Login failed:", err)
		Browser.Close()
		os.Exit(1)
	}

	exitCode := m.Run()

	Browser.Close()
	os.Exit(exitCode)
}

func Test_loginPage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Login Page",
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
			args:    args{"https://atcoder.jp/contests/acl1"},
			want:    Args{"acl1", ""},
			wantErr: false,
		},
		{
			name:    "Test #2",
			args:    args{"https://atcoder.jp/contests/m-solutions2020/tasks/m_solutions2020_a"},
			want:    Args{"m-solutions2020", "m_solutions2020_a"},
			wantErr: false,
		},
		{
			name:    "Test #3",
			args:    args{"arc107"},
			want:    Args{"arc107", ""},
			wantErr: false,
		},
		{
			name:    "Test #4", // Problem id need not match contest id.
			args:    args{"arc9999 aproblem"},
			want:    Args{"arc9999", "aproblem"},
			wantErr: false,
		},
		{
			name:    "Test #5",
			args:    args{"in_valid"},
			want:    Args{},
			wantErr: true,
		},
		{
			name:    "Test #6",
			args:    args{"in-valid in-valid"},
			want:    Args{},
			wantErr: true,
		},
		{
			name:    "Test #7",
			args:    args{""},
			want:    Args{},
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

func Test_login(t *testing.T) {
	logout()

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
			name:    "Test #1",
			args:    args{"cptools", "PleaseTryAgain"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test #2",
			args:    args{"", ""},
			want:    "",
			wantErr: true,
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

	// Hope nothing goes wrong here.
	logout()
	login(getLoginCredentials())
}
