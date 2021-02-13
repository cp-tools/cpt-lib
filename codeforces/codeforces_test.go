package codeforces

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
)

func login(usr, passwd string) (string, error) {
	p, err := loadPage(fmt.Sprintf("%v/enter", hostURL))
	if err != nil {
		return "", err
	}
	defer p.Close()

	if _, err := p.Race().Element(`#jGrowl .message`).Handle(handleErrMsg).
		Element(`#footer`).Do(); err != nil {
		return "", err
	}

	// Check if current user sesion is logged in.
	if elm := p.MustElements(`#header a[href^="/profile/"]`).First(); elm != nil {
		return elm.MustText(), nil
	}

	// Otherwise, login.
	p.MustElement("#handleOrEmail").Input(usr)
	p.MustElement("#password").Input(passwd)
	if p.MustElement("#remember").MustProperty("checked").Bool() == false {
		p.MustElement("#remember").MustClick()
	}
	p.MustElement(".submit").MustClick().WaitInvisible()

	elm, err := p.Race().Element(`.error`).Handle(handleErrMsg).
		Element(`#header a[href^="/profile/"]`).Do()
	if err != nil {
		return "", err
	}

	return elm.MustText(), nil
}

func getLoginCredentials() (string, string) {
	// setup login access to use
	usr := os.Getenv("CODEFORCES_USERNAME")
	passwd := os.Getenv("CODEFORCES_PASSWORD")
	return usr, passwd
}

func TestMain(m *testing.M) {
	// Load local .env file.
	godotenv.Load()

	_, browserHeadless := os.LookupEnv("BROWSER_HEADLESS")
	browserBin := os.Getenv("BROWSER_BINARY")
	if err := StartWithCacheDir(browserHeadless, "", browserBin, "tmp"); err != nil {
		fmt.Println("Failed to start browser:", err)
		os.Exit(1)
	}

	if handle, err := login(getLoginCredentials()); err != nil {
		fmt.Println("Login failed:", err)
		Browser.Close()
		os.Exit(1)
	} else {
		fmt.Println("Logged in user:", handle)
	}
	exitCode := m.Run()

	Browser.Close()
	os.RemoveAll("tmp")

	os.Exit(exitCode)
}

func TestArgs_String(t *testing.T) {
	tests := []struct {
		name string
		arg  Args
		want string
	}{
		{
			name: "Test #1",
			arg:  Args{"1234", "", "contest", ""},
			want: "1234 (contest)",
		},
		{
			name: "Test #2",
			arg:  Args{"1234", "b", "contest", ""},
			want: "1234 b (contest)",
		},
		{
			name: "Test #3",
			arg:  Args{"100522", "", "gym", ""},
			want: "100522 (gym)",
		},
		{
			name: "Test #4",
			arg:  Args{"100522", "f1", "gym", ""},
			want: "100522 f1 (gym)",
		},
		{
			name: "Test #5",
			arg:  Args{"201468", "", "group", "Qvv4lz52cT"},
			want: "201468 (group/Qvv4lz52cT)",
		},
		{
			name: "Test #6",
			arg:  Args{"201468", "c1", "group", "Qvv4lz52cT"},
			want: "201468 c1 (group/Qvv4lz52cT)",
		},
		{
			name: "Test #7",
			arg:  Args{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("Args.String() = %v, want %v", got, tt.want)
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
			args:    args{"https://codeforces.com/contest/1355"},
			want:    Args{"1355", "", "contest", ""},
			wantErr: false,
		},
		{
			name:    "Test #2",
			args:    args{"https://codeforces.com/contest/739/problem/E"},
			want:    Args{"739", "e", "contest", ""},
			wantErr: false,
		},
		{
			name:    "Test #3",
			args:    args{"https://codeforces.com/gym/102595"},
			want:    Args{"102595", "", "gym", ""},
			wantErr: false,
		},
		{
			name:    "Test #4",
			args:    args{"https://codeforces.com/gym/102302/problem/i"},
			want:    Args{"102302", "i", "gym", ""},
			wantErr: false,
		},
		{
			name:    "Test #5",
			args:    args{"https://codeforces.com/group/MEqF8b6wBT/contest/277493"},
			want:    Args{"277493", "", "group", "MEqF8b6wBT"},
			wantErr: false,
		},
		{
			name:    "Test #6",
			args:    args{"https://codeforces.com/group/MEqF8b6wBT/contest/277493/problem/g"},
			want:    Args{"277493", "g", "group", "MEqF8b6wBT"},
			wantErr: false,
		},
		{
			name:    "Test #7",
			args:    args{"1355"},
			want:    Args{"1355", "", "contest", ""},
			wantErr: false,
		},
		{
			name:    "Test #8",
			args:    args{"739 E"},
			want:    Args{"739", "e", "contest", ""},
			wantErr: false,
		},
		{
			name:    "Test #9",
			args:    args{"102595"},
			want:    Args{"102595", "", "gym", ""},
			wantErr: false,
		},
		{
			name:    "Test #10",
			args:    args{"102302i"},
			want:    Args{"102302", "i", "gym", ""},
			wantErr: false,
		},
		{
			name:    "Test #11",
			args:    args{"MEqF8b6wBT 277493"},
			want:    Args{"277493", "", "group", "MEqF8b6wBT"},
			wantErr: false,
		},
		{
			name:    "Test #12",
			args:    args{"MEqF8b6wBT 277493 g"},
			want:    Args{"277493", "g", "group", "MEqF8b6wBT"},
			wantErr: false,
		},
		{
			name:    "Test #13",
			args:    args{"contest"},
			want:    Args{"", "", "contest", ""},
			wantErr: false,
		},
		{
			name:    "Test #14",
			args:    args{"MEqF8b6wBT"},
			want:    Args{"", "", "group", "MEqF8b6wBT"},
			wantErr: false,
		},
		{
			name:    "Test #15",
			args:    args{"https://codeforces.com/problemset/problem/1433/E"},
			want:    Args{"1433", "e", "contest", ""},
			wantErr: false,
		},
		{
			name:    "Test #16",
			args:    args{""},
			want:    Args{},
			wantErr: false,
		},
		{
			name:    "Test #17",
			args:    args{"randomBullshitGoGo"},
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
