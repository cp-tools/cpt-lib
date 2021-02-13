package codeforces

import (
	"reflect"
	"testing"
	"time"
)

func TestArgs_GetSubmissions(t *testing.T) {
	time.Sleep(time.Second * 10)

	type args struct {
		handle string
		count  uint
	}
	tests := []struct {
		name       string
		arg        Args
		args       args
		want       []Submission
		wantErr    bool
		shouldSkip bool
	}{
		{
			name: "Test #1",
			arg:  Args{"4", "", "contest", ""},
			args: args{"cp-tools", 1},
			want: []Submission{
				{
					ID:            "81327550",
					When:          time.Date(2020, time.May, 24, 19, 14, 0, 0, time.UTC),
					Who:           "cp-tools",
					Problem:       "A - Watermelon",
					Language:      "GNU C++17",
					Verdict:       "Compilation error",
					VerdictStatus: VerdictCE,
					Time:          "0 ms",
					Memory:        "0 KB",
					IsJudging:     false,
					Arg:           Args{"4", "a", "contest", ""},
				},
				{
					ID:            "81327395",
					When:          time.Date(2020, time.May, 24, 19, 12, 0, 0, time.UTC),
					Who:           "cp-tools",
					Problem:       "A - Watermelon",
					Language:      "GNU C++17",
					Verdict:       "Compilation error",
					VerdictStatus: VerdictCE,
					Time:          "0 ms",
					Memory:        "0 KB",
					IsJudging:     false,
					Arg:           Args{"4", "a", "contest", ""},
				},
				{
					ID:            "81012854",
					When:          time.Date(2020, time.May, 23, 12, 10, 0, 0, time.UTC),
					Who:           "cp-tools",
					Problem:       "B - Before an Exam",
					Language:      "Ruby",
					Verdict:       "Runtime error on test 2",
					VerdictStatus: VerdictRTE,
					Time:          "46 ms",
					Memory:        "0 KB",
					IsJudging:     false,
					Arg:           Args{"4", "b", "contest", ""},
				},
				{
					ID:            "81011111",
					When:          time.Date(2020, time.May, 23, 11, 45, 0, 0, time.UTC),
					Who:           "cp-tools",
					Problem:       "A - Watermelon",
					Language:      "GNU C++17",
					Verdict:       "Happy New Year!",
					VerdictStatus: VerdictAC,
					Time:          "62 ms",
					Memory:        "0 KB",
					IsJudging:     false,
					Arg:           Args{"4", "a", "contest", ""},
				},
			},
			wantErr: false,
		},
		{
			name: "Test #2",
			arg:  Args{"4", "b", "contest", ""},
			args: args{"cp-tools", 1e9},
			want: []Submission{
				{
					ID:            "81012854",
					When:          time.Date(2020, time.May, 23, 12, 10, 0, 0, time.UTC),
					Who:           "cp-tools",
					Problem:       "B - Before an Exam",
					Language:      "Ruby",
					Verdict:       "Runtime error on test 2",
					VerdictStatus: VerdictRTE,
					Time:          "46 ms",
					Memory:        "0 KB",
					IsJudging:     false,
					Arg:           Args{"4", "b", "contest", ""},
				},
			},
			wantErr: false,
		},
		{
			name:       "Test #3",
			arg:        Args{"5", "", "contest", ""},
			args:       args{"cp-tools", 2},
			want:       nil, // Being skipped.
			wantErr:    false,
			shouldSkip: true,
		},
		{
			name:    "Test #4",
			arg:     Args{"207982", "", "group", "7rY4CfQSjd"},
			args:    args{"invalid", 1},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test #5",
			arg:     Args{"12345", "", "contest", ""},
			args:    args{"cp-tools", 1},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.arg.GetSubmissions(tt.args.handle, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.GetSubmissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				// No data is returned; continue.
				return
			}

			// read till channel closes
			submissions := make([]Submission, 0)
			for v := range got {
				t.Log("Data rows in page:", len(v))
				submissions = append(submissions, v...)
			}

			if tt.shouldSkip {
				// Check for duplicates.
				tmpMap := make(map[Submission]bool)
				for _, submission := range submissions {
					tmpMap[submission] = true
				}

				if len(tmpMap) != len(submissions) {
					t.Errorf("Args.GetSubmissions() returned duplicate values")
				}
				if uint(len(submissions)) != 100 { // There should be exactly this many.
					t.Errorf("Args.GetSubmissions() required = 100 rows, got %v rows", len(submissions))
				}
				// All fine.
				t.SkipNow()
			}

			if !reflect.DeepEqual(submissions, tt.want) {
				t.Errorf("Args.GetSubmissions() = %v, want %v", submissions, tt.want)
			}
		})
	}
}

func TestSubmission_GetSourceCode(t *testing.T) {
	time.Sleep(time.Second * 10)

	tests := []struct {
		name    string
		sub     Submission
		want    string
		wantErr bool
	}{
		{
			name: "Test #1",
			sub:  Submission{ID: "81012854", Arg: Args{"4", "", "contest", ""}}, // just bare info here
			want: `d,s=gets.split.map(&:to_i)
d.times{$*<<gets.split.map(&:to_i)}
a=$*.transpose
x=s-a[0].inject(:+)
puts x<0||s>a[1].inject(:+) ?:NO:"YES
"+$*.map{|l,r|t=[r-l,x].min;x-=t;l+t}*" "`,
			wantErr: false,
		},
		{
			name: "Test #2",
			sub:  Submission{ID: "95913201", Arg: Args{"1359", "", "contest", ""}}, // just bare info here
			want: `//go corona go
#include <bits/stdc++.h>
#include<cmath>
using namespace std;

int main() {
	int n;
	cin>>n;
	vector<int> arr(n);
	for(int i=0; i<n; i++){
		cin>>arr[i];
	}
	vector<int> dp(n);
	vector<int> maxi(n);
	dp[0]=arr[0];
	maxi[0]=dp[0];
	int m=0;
	for(int i=1; i<n; i++){
		if((dp[i-1]+arr[i])>arr[i]){
			dp[i]=dp[i-1]+arr[i];
			maxi[i]=max(maxi[i-1],arr[i]);
		}
		else{
			dp[i]=arr[i];
			maxi[i]=arr[i];
		}
		m=max(m,(dp[i]-maxi[i]));
	}
	cout<<m<<endl;
}`,
			wantErr: false,
		},
		{
			name:    "Test #3",
			sub:     Submission{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test #4",
			sub:     Submission{ID: "12345678", Arg: Args{"4", "b", "contest", ""}},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sub.GetSourceCode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Submission.GetSourceCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Submission.GetSourceCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
