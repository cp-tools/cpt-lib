package atcoder

import "testing"

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

func TestArgs_DashboardPage(t *testing.T) {
	tests := []struct {
		name     string
		arg      Args
		wantLink string
		wantErr  bool
	}{
		{
			name:     "Test #1",
			arg:      Args{"hhkb2020", ""},
			wantLink: "https://atcoder.jp/contests/hhkb2020",
			wantErr:  false,
		},
		{
			name:     "Test #2",
			arg:      Args{},
			wantLink: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLink, err := tt.arg.DashboardPage()
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.DashboardPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLink != tt.wantLink {
				t.Errorf("Args.DashboardPage() = %v, want %v", gotLink, tt.wantLink)
			}
		})
	}
}

func TestArgs_VirtualPage(t *testing.T) {
	tests := []struct {
		name     string
		arg      Args
		wantLink string
		wantErr  bool
	}{
		{
			name:     "Test #1",
			arg:      Args{"tokiomarine2020", ""},
			wantLink: "https://atcoder.jp/contests/tokiomarine2020/virtual",
			wantErr:  false,
		},
		{
			name:     "Test #2",
			arg:      Args{},
			wantLink: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLink, err := tt.arg.VirtualPage()
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.DashboardPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLink != tt.wantLink {
				t.Errorf("Args.DashboardPage() = %v, want %v", gotLink, tt.wantLink)
			}
		})
	}
}
