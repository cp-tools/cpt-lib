package codeforces

import (
	"reflect"
	"testing"
)

func TestArgs_countdownPage(t *testing.T) {
	tests := []struct {
		name    string
		arg     Args
		want    string
		wantErr bool
	}{
		{
			name:    "Test #1",
			arg:     Args{"1234", "", "contest", ""},
			want:    "https://codeforces.com/contest/1234/countdown",
			wantErr: false,
		},
		{
			name:    "Test #2",
			arg:     Args{"100001", "", "gym", ""},
			want:    "https://codeforces.com/gym/100001/countdown",
			wantErr: false,
		},
		{
			name:    "Test #3",
			arg:     Args{"277493", "", "group", "MEqF8b6wBT"},
			want:    "https://codeforces.com/group/MEqF8b6wBT/contest/277493/countdown",
			wantErr: false,
		},
		{
			name:    "Test #4",
			arg:     Args{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test #5",
			arg:     Args{"288493", "", "group", ""},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.CountdownPage()
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.countdownPage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Args.countdownPage() = %v, want %v", got, tt.want)

			}
		})
	}
}

func TestArgs_contestsPage(t *testing.T) {
	tests := []struct {
		name    string
		arg     Args
		want    string
		wantErr bool
	}{
		{
			name:    "Test #1",
			arg:     Args{"1234", "", "contest", ""},
			want:    "https://codeforces.com/contests/1234",
			wantErr: false,
		},
		{
			name:    "Test #2",
			arg:     Args{"100001", "", "gym", ""},
			want:    "https://codeforces.com/contests/100001",
			wantErr: false,
		},
		{
			name:    "Test #3",
			arg:     Args{"277493", "", "group", "MEqF8b6wBT"},
			want:    "https://codeforces.com/group/MEqF8b6wBT/contests?complete=true",
			wantErr: false,
		},
		{
			name:    "Test #4",
			arg:     Args{"", "", "contest", ""},
			want:    "https://codeforces.com/contests?complete=true",
			wantErr: false,
		},
		{
			name:    "Test #5",
			arg:     Args{"", "", "gym", ""},
			want:    "https://codeforces.com/gyms?complete=true",
			wantErr: false,
		},
		{
			name:    "Test #6",
			arg:     Args{"288493", "", "group", ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test #7",
			arg:     Args{},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.ContestsPage()
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.contestsPage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Args.contestsPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgs_dashboardPage(t *testing.T) {
	tests := []struct {
		name    string
		arg     Args
		want    string
		wantErr bool
	}{
		{
			name:    "Test #1",
			arg:     Args{"1234", "", "contest", ""},
			want:    "https://codeforces.com/contest/1234",
			wantErr: false,
		},
		{
			name:    "Test #2",
			arg:     Args{"100001", "", "gym", ""},
			want:    "https://codeforces.com/gym/100001",
			wantErr: false,
		},
		{
			name:    "Test #3",
			arg:     Args{"277493", "", "group", "MEqF8b6wBT"},
			want:    "https://codeforces.com/group/MEqF8b6wBT/contest/277493",
			wantErr: false,
		},
		{
			name:    "Test #4",
			arg:     Args{"288493", "", "group", ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test #5",
			arg:     Args{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test #6",
			arg:     Args{"123", "", "invalid", ""},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.DashboardPage()
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.dashboardPage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Args.dashboardPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgs_registerPage(t *testing.T) {
	tests := []struct {
		name    string
		arg     Args
		want    string
		wantErr bool
	}{
		{
			name:    "Test #1",
			arg:     Args{"1234", "", "contest", ""},
			want:    "https://codeforces.com/contestRegistration/1234",
			wantErr: false,
		},
		{
			name:    "Test #2",
			arg:     Args{"", "", "contest", ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test #3",
			arg:     Args{"288493", "", "", ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test #4",
			arg:     Args{"288493", "", "invalid", ""},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.RegisterPage()
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.registerPage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Args.registerPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgs_problemsPage(t *testing.T) {
	tests := []struct {
		name    string
		arg     Args
		want    string
		wantErr bool
	}{
		{
			name:    "Test #1",
			arg:     Args{"1", "", "contest", ""},
			want:    "https://codeforces.com/contest/1/problems",
			wantErr: false,
		},
		{
			name:    "Test #2",
			arg:     Args{"4", "b", "contest", ""},
			want:    "https://codeforces.com/contest/4/problem/b",
			wantErr: false,
		},
		{
			name:    "Test #3",
			arg:     Args{"102341", "", "gym", ""},
			want:    "https://codeforces.com/gym/102341/problems",
			wantErr: false,
		},
		{
			name:    "Test #4",
			arg:     Args{"102323", "a", "gym", ""},
			want:    "https://codeforces.com/gym/102323/problem/a",
			wantErr: false,
		},
		{
			name:    "Test #5",
			arg:     Args{"283855", "", "group", "bK73bvp3d7"},
			want:    "https://codeforces.com/group/bK73bvp3d7/contest/283855/problems",
			wantErr: false,
		},
		{
			name:    "Test #6",
			arg:     Args{"283855", "c", "group", "bK73bvp3d7"},
			want:    "https://codeforces.com/group/bK73bvp3d7/contest/283855/problem/c",
			wantErr: false,
		},
		{
			name:    "Test #7",
			arg:     Args{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test #8",
			arg:     Args{"283855", "", "group", ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test #9",
			arg:     Args{"45", "d", "invalid", ""},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.ProblemsPage()
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.problemsPage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Args.problemsPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgs_submissionsPage(t *testing.T) {
	type args struct {
		handle string
	}
	tests := []struct {
		name    string
		arg     Args
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Test #1",
			arg:     Args{},
			args:    args{"cp-tools"},
			want:    "https://codeforces.com/submissions/cp-tools",
			wantErr: false,
		},
		{
			name:    "Test #2",
			arg:     Args{"4", "a", "contest", ""},
			args:    args{"cp-tools"},
			want:    "https://codeforces.com/submissions/cp-tools/contest/4",
			wantErr: false,
		},
		{
			name:    "Test #3",
			arg:     Args{"102595", "", "gym", ""},
			args:    args{"cp-tools"},
			want:    "https://codeforces.com/submissions/cp-tools/gym/102595",
			wantErr: false,
		},
		{
			name:    "Test #4",
			arg:     Args{},
			args:    args{""},
			want:    "https://codeforces.com/submissions/cp-tools",
			wantErr: false,
		},
		{
			name:    "Test #5",
			arg:     Args{"207982", "", "group", "7rY4CfQSjd"},
			args:    args{"invalid"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test #6",
			arg:     Args{"207982", "", "group", "7rY4CfQSjd"},
			args:    args{""},
			want:    "https://codeforces.com/group/7rY4CfQSjd/contest/207982/my",
			wantErr: false,
		},
		{
			name:    "Test #7",
			arg:     Args{"965", "", "contest", ""},
			args:    args{""},
			want:    "https://codeforces.com/contest/965/my",
			wantErr: false,
		},
		{
			name:    "Test #8",
			arg:     Args{"50", "", "randomBullshitGoGo", ""},
			args:    args{"cp-tools"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Test #4" {
				t.SkipNow()
			}

			got, err := tt.arg.SubmissionsPage(tt.args.handle)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.submissionsPage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Args.submissionsPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubmission_sourceCodePage(t *testing.T) {
	tests := []struct {
		name    string
		sub     Submission
		want    string
		wantErr bool
	}{
		{
			name:    "Test #1",
			sub:     Submission{ID: "81011111", Arg: Args{"4", "", "contest", ""}}, // rest is not required
			want:    "https://codeforces.com/contest/4/submission/81011111",
			wantErr: false,
		},
		{
			name:    "Test #2",
			sub:     Submission{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test #3",
			sub:     Submission{ID: "95913201", Arg: Args{"207982", "", "group", "7rY4CfQSjd"}},
			want:    "https://codeforces.com/group/7rY4CfQSjd/contest/207982/submission/95913201",
			wantErr: false,
		},
		{
			name:    "Test #4",
			sub:     Submission{ID: "1234567", Arg: Args{"4", "", "invalid", ""}},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sub.SourceCodePage()
			if (err != nil) != tt.wantErr {
				t.Errorf("Submission.sourceCodePage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Submission.sourceCodePage() = %v, want %v", got, tt.want)
			}
		})
	}
}
