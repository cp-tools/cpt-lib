package atcoder

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/go-rod/rod"
	"github.com/joho/godotenv"
)

func login(usr, passwd string) (string, error) {
	p, err := loadPage(fmt.Sprintf("%v/login", hostURL))
	if err != nil {
		return "", err
	}
	defer p.Close()

	if _, err := p.Race().Element(`alert`).Handle(handleErrMsg).
		Element(`footer.footer`).Do(); err != nil {
		return "", err
	}

	// Check if current user is logged in.
	if handle := p.MustEval(`userScreenName`).String(); handle != "" {
		return handle, nil
	}
	// Otherwise, login.
	p.MustElement("#username").Input(usr)
	p.MustElement("#password").Input(passwd)
	p.MustElement("#submit").MustClick().WaitInvisible()

	if _, err := p.Race().ElementR(`.alert`, `Username or Password is incorrect`).
		Handle(func(e *rod.Element) error { return fmt.Errorf(e.MustText()) }).
		Element(`.navbar-right>li:last-child>a[class]`).Do(); err != nil {
		return "", err
	}

	handle := p.MustEval(`userScreenName`).String()
	return handle, nil
}

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
