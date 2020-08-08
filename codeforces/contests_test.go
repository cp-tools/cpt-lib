package codeforces

import (
	"reflect"
	"testing"
	"time"
)

func TestArgs_countdownPage(t *testing.T) {
	tests := []struct {
		name     string
		arg      Args
		wantLink string
	}{
		{
			name:     "Test #1",
			arg:      Args{"1234", "", "contest", ""},
			wantLink: "https://codeforces.com/contest/1234/countdown",
		},
		{
			name:     "Test #2",
			arg:      Args{"100001", "", "gym", ""},
			wantLink: "https://codeforces.com/gym/100001/countdown",
		},
		{
			name:     "Test #3",
			arg:      Args{"277493", "", "group", "MEqF8b6wBT"},
			wantLink: "https://codeforces.com/group/MEqF8b6wBT/contest/277493/countdown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLink := tt.arg.CountdownPage(); gotLink != tt.wantLink {
				t.Errorf("Args.countdownPage() = %v, want %v", gotLink, tt.wantLink)
			}
		})
	}
}

func TestArgs_contestsPage(t *testing.T) {
	tests := []struct {
		name     string
		arg      Args
		wantLink string
	}{
		{
			name:     "Test #1",
			arg:      Args{"1234", "", "contest", ""},
			wantLink: "https://codeforces.com/contests/1234",
		},
		{
			name:     "Test #2",
			arg:      Args{"100001", "", "gym", ""},
			wantLink: "https://codeforces.com/contests/100001",
		},
		{
			name:     "Test #3",
			arg:      Args{"277493", "", "group", "MEqF8b6wBT"},
			wantLink: "https://codeforces.com/group/MEqF8b6wBT/contests?complete=true",
		},
		{
			name:     "Test #4",
			arg:      Args{"", "", "contest", ""},
			wantLink: "https://codeforces.com/contests?complete=true",
		},
		{
			name:     "Test #5",
			arg:      Args{"", "", "gym", ""},
			wantLink: "https://codeforces.com/gyms?complete=true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLink := tt.arg.ContestsPage(); gotLink != tt.wantLink {
				t.Errorf("Args.contestsPage() = %v, want %v", gotLink, tt.wantLink)
			}
		})
	}
}

func TestArgs_registerPage(t *testing.T) {
	tests := []struct {
		name     string
		arg      Args
		wantLink string
	}{
		{
			name:     "Test #1",
			arg:      Args{"1234", "", "contest", ""},
			wantLink: "https://codeforces.com/contestRegistration/1234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLink := tt.arg.RegisterPage(); gotLink != tt.wantLink {
				t.Errorf("Args.registerPage() = %v, want %v", gotLink, tt.wantLink)
			}
		})
	}
}

func TestArgs_GetCountdown(t *testing.T) {
	tests := []struct {
		name    string
		arg     Args
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "Test #1",
			arg:     Args{"1234", "", "contest", ""},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Test #2",
			arg:     Args{"100001", "", "gym", ""},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Test #3",
			arg:     Args{"283855", "", "group", "bK73bvp3d7"},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Test #4",
			arg:     Args{"12345", "", "contest", ""},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.GetCountdown()
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.GetCountdown() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Args.GetCountdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgs_GetContests(t *testing.T) {
	type args struct {
		omitFinishedContests bool
	}
	tests := []struct {
		name    string
		arg     Args
		args    args
		want    []Contest
		wantErr bool
	}{
		{
			name: "Test #1",
			arg:  Args{"7", "", "contest", ""},
			args: args{false},
			want: []Contest{
				{
					Name:        "Codeforces Beta Round #7",
					Writers:     []string{"MikeMirzayanov", "RAD", "e-maxx"},
					StartTime:   time.Date(2010, time.April, 1, 16, 45, 0, 0, time.UTC),
					Duration:    time.Hour * 2,
					RegCount:    722,
					RegStatus:   RegistrationClosed,
					Description: nil,
					Arg:         Args{"7", "", "contest", ""},
				},
			},
			wantErr: false,
		},
		{
			name: "Test #2",
			arg:  Args{"100499", "", "gym", ""},
			args: args{false},
			want: []Contest{
				{
					Name:        "2014 ACM-ICPC Vietnam National First Round",
					Writers:     nil,
					StartTime:   time.Date(2014, time.October, 12, 7, 0, 0, 0, time.UTC),
					Duration:    time.Hour*5 + time.Minute*15,
					RegCount:    RegistrationNotExists,
					RegStatus:   RegistrationNotExists,
					Description: []string{"Prepared by I_love_Hoang_Yen"},
					Arg:         Args{"100499", "", "gym", ""},
				},
			},
			wantErr: false,
		},
		{
			name: "Test #3",
			arg:  Args{"", "", "group", "7rY4CfQSjd"},
			args: args{false},
			want: []Contest{
				{
					Name:        "gym problems -2",
					Writers:     nil,
					StartTime:   time.Date(2016, time.July, 19, 6, 30, 0, 0, time.UTC),
					Duration:    time.Hour * 4,
					RegCount:    RegistrationNotExists,
					RegStatus:   RegistrationNotExists,
					Description: []string{"Prepared by Daniar", "Training Camp Contest", "Syria, Homs", "Statements:\nin English"},
					Arg:         Args{"207982", "", "group", "7rY4CfQSjd"},
				},
				{
					Name:        "gym problems",
					Writers:     nil,
					StartTime:   time.Date(2016, time.July, 18, 7, 0, 0, 0, time.UTC),
					Duration:    time.Hour * 4,
					RegCount:    RegistrationNotExists,
					RegStatus:   RegistrationNotExists,
					Description: []string{"Prepared by Daniar", "Training Camp Contest", "Syria, Homs"},
					Arg:         Args{"207960", "", "group", "7rY4CfQSjd"},
				},
				{
					Name:        "Al-Baath Training Camp 2016 - Advanced Contest",
					Writers:     nil,
					StartTime:   time.Date(2016, time.March, 12, 8, 30, 0, 0, time.UTC),
					Duration:    time.Hour * 3,
					RegCount:    RegistrationNotExists,
					RegStatus:   RegistrationNotExists,
					Description: []string{"Prepared by sqr_hussain", "Training Camp Contest", "Syria, Glenroy, 2016-2017"},
					Arg:         Args{"206484", "", "group", "7rY4CfQSjd"},
				},
				{
					Name:      "Al-Baath Training Camp 2016 - Beginners Contest",
					Writers:   nil,
					StartTime: time.Date(2016, time.March, 12, 8, 30, 0, 0, time.UTC),
					Duration:  time.Hour * 3,
					RegCount:  RegistrationNotExists,
					RegStatus: RegistrationNotExists,
					Description: []string{"Prepared by sqr_hussain", "Official International Personal Contest",
						"Syria, Glenroy, 2016-2017"},
					Arg: Args{"206482", "", "group", "7rY4CfQSjd"},
				},
				{
					Name:        "Al-Baath Training day-10 Beginners",
					Writers:     nil,
					StartTime:   time.Date(2016, time.March, 2, 7, 30, 0, 0, time.UTC),
					Duration:    time.Hour*2 + time.Minute*30,
					RegCount:    RegistrationNotExists,
					RegStatus:   RegistrationNotExists,
					Description: []string{"Prepared by Daniar"},
					Arg:         Args{"206359", "", "group", "7rY4CfQSjd"},
				},
				{
					Name:        "ALBAATH Rush day 9 Intermediate",
					Writers:     nil,
					StartTime:   time.Unix(0, 0).UTC(),
					Duration:    time.Hour*2 + time.Minute*30,
					RegCount:    RegistrationNotExists,
					RegStatus:   RegistrationNotExists,
					Description: []string{"Prepared by Marcil"},
					Arg:         Args{"206346", "", "group", "7rY4CfQSjd"},
				},
				{
					Name:        "ALBAATH Rush day 9 Begginners",
					Writers:     nil,
					StartTime:   time.Unix(0, 0).UTC(),
					Duration:    time.Hour*2 + time.Minute*30,
					RegCount:    RegistrationNotExists,
					RegStatus:   RegistrationNotExists,
					Description: []string{"Prepared by Marcil"},
					Arg:         Args{"206344", "", "group", "7rY4CfQSjd"},
				},
			},
			wantErr: false,
		},
		{
			name:    "Test #4",
			arg:     Args{"7", "", "contest", ""},
			args:    args{true},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.GetContests(tt.args.omitFinishedContests)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.GetContests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Args.GetContests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgs_RegisterForContest(t *testing.T) {
	tests := []struct {
		name    string
		arg     Args
		want    *RegisterInfo
		wantErr bool
	}{
		// TODO: Implement tests for RegisterForContest()
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.RegisterForContest()
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.RegisterForContest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Args.RegisterForContest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgs_dashboardPage(t *testing.T) {
	tests := []struct {
		name     string
		arg      Args
		wantLink string
	}{
		{
			name:     "Test #1",
			arg:      Args{"1234", "", "contest", ""},
			wantLink: "https://codeforces.com/contest/1234",
		},
		{
			name:     "Test #2",
			arg:      Args{"100001", "", "gym", ""},
			wantLink: "https://codeforces.com/gym/100001",
		},
		{
			name:     "Test #3",
			arg:      Args{"277493", "", "group", "MEqF8b6wBT"},
			wantLink: "https://codeforces.com/group/MEqF8b6wBT/contest/277493",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLink := tt.arg.DashboardPage(); gotLink != tt.wantLink {
				t.Errorf("Args.dashboardPage() = %v, want %v", gotLink, tt.wantLink)
			}
		})
	}
}

func TestArgs_GetDashboard(t *testing.T) {
	tests := []struct {
		name    string
		arg     Args
		want    Dashboard
		wantErr bool
	}{
		{
			name: "Test #1",
			arg:  Args{"4", "", "contest", ""},
			want: Dashboard{
				Name: "Codeforces Beta Round #4 (Div. 2 Only)",
				Problem: []Problem{
					{
						Name:        "Watermelon",
						TimeLimit:   "1 s",
						MemoryLimit: "64 MB",
						InpStream:   "standard input",
						OutStream:   "standard output",
						SampleTests: nil,
						SolveCount:  -1, // keeps changing, ignore value
						SolveStatus: SolveAccepted,
						Arg:         Args{"4", "a", "contest", ""},
					},
					{
						Name:        "Before an Exam",
						TimeLimit:   "0.5 s",
						MemoryLimit: "64 MB",
						InpStream:   "standard input",
						OutStream:   "standard output",
						SampleTests: nil,
						SolveCount:  -1, // keeps changing, ignore value
						SolveStatus: SolveRejected,
						Arg:         Args{"4", "b", "contest", ""},
					},
					{
						Name:        "Registration System",
						TimeLimit:   "5 s",
						MemoryLimit: "64 MB",
						InpStream:   "standard input",
						OutStream:   "standard output",
						SampleTests: nil,
						SolveCount:  -1, // keeps changing, ignore value
						SolveStatus: SolveNotAttempted,
						Arg:         Args{"4", "c", "contest", ""},
					},
					{
						Name:        "Mysterious Present",
						TimeLimit:   "1 s",
						MemoryLimit: "64 MB",
						InpStream:   "standard input",
						OutStream:   "standard output",
						SampleTests: nil,
						SolveCount:  -1, // keeps changing, ignore value
						SolveStatus: SolveNotAttempted,
						Arg:         Args{"4", "d", "contest", ""},
					},
				},
				Countdown: 0,
				Material: map[string]string{
					"https://codeforces.com/blog/entry/158": "Announcement",
					"https://codeforces.com/blog/entry/161": "Tutorial #1 (en)",
					"https://codeforces.com/blog/entry/163": "Tutorial #2 (en)",
					"https://codeforces.com/blog/entry/164": "Tutorial #3 (ru)",
					"https://codeforces.com/blog/entry/178": "Tutorial #4",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.GetDashboard()
			// set solve count to -1
			for i := range got.Problem {
				got.Problem[i].SolveCount = -1
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Args.GetDashboard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Args.GetDashboard() = %v, want %v", got, tt.want)
			}
		})
	}
}
