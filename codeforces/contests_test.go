package codeforces

import (
	"reflect"
	"testing"
	"time"
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
			wantErr: true, // No such contest
		},
		{
			name:    "Test #5",
			arg:     Args{"942", "", "", ""},
			want:    0,
			wantErr: true, // You are not allowed to view the contest
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
		pageCount uint
	}
	tests := []struct {
		name       string
		arg        Args
		args       args
		want       []Contest
		wantErr    bool
		shouldSkip bool
	}{
		{
			name: "Test #1",
			arg:  Args{"7", "", "contest", ""},
			args: args{1},
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
			args: args{1e9},
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
			args: args{2},
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
			name:       "Test #4",
			arg:        Args{"", "", "gym", ""},
			args:       args{4},
			want:       nil, // Being skipped.
			wantErr:    false,
			shouldSkip: true,
		},
		{
			name:       "Test #5",
			arg:        Args{"", "", "contest", ""},
			args:       args{2},
			want:       nil, // Being skipped.
			wantErr:    false,
			shouldSkip: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.GetContests(tt.args.pageCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.GetContests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			contests := make([]Contest, 0)
			for v := range got {
				contests = append(contests, v...)
			}

			if tt.shouldSkip {
				// Check if there are duplicates.
				tmpMap := make(map[Args]bool)
				for _, contest := range contests {
					tmpMap[contest.Arg] = true
				}

				if len(tmpMap) != len(contests) {
					t.Errorf("Args.GetContests() returned duplicate values")
				}
				if uint(len(contests)) < 100*tt.args.pageCount {
					t.Errorf("Args.GetContests() required >= %v rows, got %v rows",
						100*tt.args.pageCount, len(contests))
				}
				// No duplicates found.
				t.SkipNow()
			}

			if !reflect.DeepEqual(contests, tt.want) {
				t.Errorf("Args.GetContests() = %v, want %v", contests, tt.want)
			}
		})
	}
}

func TestArgs_RegisterForContest(t *testing.T) {
	// check if there exists any contests
	// with open registration available
	arg := Args{Class: ClassContest}
	chanContests, err := arg.GetContests(1)
	if err != nil {
		t.Errorf("arg.GetContests() error = %v", err)
	}

	for contests := range chanContests {
		for _, cont := range contests {
			if cont.RegStatus != RegistrationOpen {
				continue
			}

			// found contest. Try registering in it.
			regInfo, err := cont.Arg.RegisterForContest()
			if err != nil {
				t.Errorf("cont.Arg.RegisterForContest() error = %v", err)
			}

			err = regInfo.Register()
			if err != nil {
				t.Errorf("regInfo.Register() error = %v", err)
			}

			// check registered status
			chanCurrContests, err := cont.Arg.GetContests(1)
			if err != nil {
				t.Errorf("cont.Arg.GetContests() error = %v", err)
			}

			for currContests := range chanCurrContests {
				// check if registered in contest
				if currContests[0].RegStatus != RegistrationDone {
					t.Errorf("Registration failed - %v", currContests[0])
				}
				t.Logf("Successfully registered in %v", cont)
			}

			break
		}
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
				Material:  map[string]string{
					// some error caused them to disappear
					// format (for reference) is link->title
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
