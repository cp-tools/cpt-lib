package codeforces

import (
	"reflect"
	"testing"
	"time"
)

// @todo Add more tests in submissions file

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestArgs_GetSubmissions(t *testing.T) {
	type args struct {
		handle string
		count  int
	}
	tests := []struct {
		name    string
		arg     Args
		args    args
		want    []Submission
		wantErr bool
	}{
		{
			name: "Test #1",
			arg:  Args{"4", "", "contest", ""},
			args: args{"cp-tools", 1},
			want: []Submission{
				{
					ID:        "81327550",
					When:      time.Date(2020, time.May, 24, 19, 14, 0, 0, time.UTC),
					Who:       "cp-tools",
					Problem:   "A - Watermelon",
					Language:  "GNU C++17",
					Verdict:   "Compilation error",
					Time:      "0 ms",
					Memory:    "0 KB",
					IsJudging: false,
					Arg:       Args{"4", "a", "contest", ""},
				},
				{
					ID:        "81327395",
					When:      time.Date(2020, time.May, 24, 19, 12, 0, 0, time.UTC),
					Who:       "cp-tools",
					Problem:   "A - Watermelon",
					Language:  "GNU C++17",
					Verdict:   "Compilation error",
					Time:      "0 ms",
					Memory:    "0 KB",
					IsJudging: false,
					Arg:       Args{"4", "a", "contest", ""},
				},
				{
					ID:        "81012854",
					When:      time.Date(2020, time.May, 23, 12, 10, 0, 0, time.UTC),
					Who:       "cp-tools",
					Problem:   "B - Before an Exam",
					Language:  "Ruby",
					Verdict:   "Runtime error on test 2",
					Time:      "46 ms",
					Memory:    "0 KB",
					IsJudging: false,
					Arg:       Args{"4", "b", "contest", ""},
				},
				{
					ID:        "81011111",
					When:      time.Date(2020, time.May, 23, 11, 45, 0, 0, time.UTC),
					Who:       "cp-tools",
					Problem:   "A - Watermelon",
					Language:  "GNU C++17",
					Verdict:   "Accepted",
					Time:      "62 ms",
					Memory:    "0 KB",
					IsJudging: false,
					Arg:       Args{"4", "a", "contest", ""},
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
					ID:        "81012854",
					When:      time.Date(2020, time.May, 23, 12, 10, 0, 0, time.UTC),
					Who:       "cp-tools",
					Problem:   "B - Before an Exam",
					Language:  "Ruby",
					Verdict:   "Runtime error on test 2",
					Time:      "46 ms",
					Memory:    "0 KB",
					IsJudging: false,
					Arg:       Args{"4", "b", "contest", ""},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.arg.GetSubmissions(tt.args.handle, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.GetSubmissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// read till channel closes
			val := make([]Submission, 0)
			for v := range got {
				val = append(val, v...)
			}

			if !reflect.DeepEqual(val, tt.want) {
				t.Errorf("Args.GetSubmissions() = %v, want %v", val, tt.want)
			}
		})
	}
}

/*
func TestArgs_ChanGetSubmissions(t *testing.T) {
	// get (some) submission to submit
	arg, _ := Parse("1132c")
	chanSubmissions, _ := arg.GetSubmissions("hohomu", 1)
	submissions := <-chanSubmissions

	sourceCode, _ := submissions[0].GetSourceCode()
	sourceCode += fmt.Sprintf("\n//%v\n", genRandomString(10))
	file, _ := ioutil.TempFile(os.TempDir(), "tmp-source")
	defer os.Remove(file.Name())
	file.WriteString(sourceCode)

	// submit solution, and monitor submission status
	arg.SubmitSolution("GNU G++14 6.4.0", file.Name())
	chanSubmissions, _ = arg.GetSubmissions("cp-tools", 1)
	for submissions := range chanSubmissions {
		t.Log(submissions[0].Verdict)
	}
	// works as expected!
}
*/

func TestSubmission_GetSourceCode(t *testing.T) {
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
