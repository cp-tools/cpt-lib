package codeforces

import (
	"reflect"
	"testing"
)

func TestArgs_problemsPage(t *testing.T) {
	tests := []struct {
		name     string
		arg      Args
		wantLink string
	}{
		{
			name:     "Test #1",
			arg:      Args{"1", "", "contest", ""},
			wantLink: "https://codeforces.com/contest/1/problems",
		},
		{
			name:     "Test #2",
			arg:      Args{"4", "b", "contest", ""},
			wantLink: "https://codeforces.com/contest/4/problem/b",
		},
		{
			name:     "Test #3",
			arg:      Args{"102341", "", "gym", ""},
			wantLink: "https://codeforces.com/gym/102341/problems",
		},
		{
			name:     "Test #4",
			arg:      Args{"102323", "a", "gym", ""},
			wantLink: "https://codeforces.com/gym/102323/problem/a",
		},
		{
			name:     "Test #5",
			arg:      Args{"283855", "", "group", "bK73bvp3d7"},
			wantLink: "https://codeforces.com/group/bK73bvp3d7/contest/283855/problems",
		},
		{
			name:     "Test #6",
			arg:      Args{"283855", "c", "group", "bK73bvp3d7"},
			wantLink: "https://codeforces.com/group/bK73bvp3d7/contest/283855/problem/c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLink := tt.arg.problemsPage(); gotLink != tt.wantLink {
				t.Errorf("Args.problemsPage() = %v, want %v", gotLink, tt.wantLink)
			}
		})
	}
}

func TestArgs_GetProblems(t *testing.T) {
	tests := []struct {
		name    string
		arg     Args
		want    []Problem
		wantErr bool
	}{
		{
			name: "Test #1",
			arg:  Args{"4", "", "contest", ""},
			want: []Problem{
				{
					Name:        "A. Watermelon",
					TimeLimit:   "1 second",
					MemoryLimit: "64 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "8\n",
							Output: "YES\n",
						},
					},
					Arg: Args{"4", "a", "contest", ""},
				},
				{
					Name:        "B. Before an Exam",
					TimeLimit:   "0.5 second",
					MemoryLimit: "64 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "1 48\n5 7\n",
							Output: "NO\n",
						},
						{
							Input:  "2 5\n0 1\n3 5\n",
							Output: "YES\n1 4\n",
						},
					},
					Arg: Args{"4", "b", "contest", ""},
				},
				{
					Name:        "C. Registration system",
					TimeLimit:   "5 seconds",
					MemoryLimit: "64 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "4\nabacaba\nacaba\nabacaba\nacab\n",
							Output: "OK\nOK\nabacaba1\nOK\n",
						},
						{
							Input:  "6\nfirst\nfirst\nsecond\nsecond\nthird\nthird\n",
							Output: "OK\nfirst1\nOK\nsecond1\nOK\nthird1\n",
						},
					},
					Arg: Args{"4", "c", "contest", ""},
				},
				{
					Name:        "D. Mysterious Present",
					TimeLimit:   "1 second",
					MemoryLimit: "64 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "2 1 1\n2 2\n2 2\n",
							Output: "1\n1\n",
						},
						{
							Input:  "3 3 3\n5 4\n12 11\n9 8\n",
							Output: "3\n1 3 2\n",
						},
					},
					Arg: Args{"4", "d", "contest", ""},
				},
			},
			wantErr: false,
		},
		{
			name: "Test #2",
			arg:  Args{"1234", "a", "contest", ""},
			want: []Problem{
				{
					Name:        "A. Equalize Prices Again",
					TimeLimit:   "1 second",
					MemoryLimit: "256 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "3\n5\n1 2 3 4 5\n3\n1 2 2\n4\n1 1 1 1\n",
							Output: "3\n2\n1\n",
						},
					},
					Arg: Args{"1234", "a", "contest", ""},
				},
			},
			wantErr: false,
		},
		{
			name: "Test #3",
			arg:  Args{"101189", "", "gym", ""},
			want: []Problem{
				{
					Name:        "A. Arpa’s hard exam and Mehrdad’s naive cheat(Hard)",
					TimeLimit:   "1 second",
					MemoryLimit: "256 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "1\n",
							Output: "8\n",
						},
						{
							Input:  "2\n",
							Output: "4\n",
						},
					},
					Arg: Args{"101189", "a", "gym", ""},
				},
				{
					Name:        "B. Arpa’s obvious problem and Mehrdad’s terrible solution(Hard)",
					TimeLimit:   "0.5 seconds",
					MemoryLimit: "512 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "2 3\n1 2\n",
							Output: "1\n",
						},
						{
							Input:  "4 1\nA C E F\n",
							Output: "1\n",
						},
					},
					Arg: Args{"101189", "b", "gym", ""},
				},
				{
					Name:        "C. Arpa's loud Owf and Mehrdad's evil plan(Hard)",
					TimeLimit:   "3 seconds",
					MemoryLimit: "256 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "4\n2 3 1 4\n",
							Output: "3\n",
						},
						{
							Input:  "4\n4 4 4 4\n",
							Output: "-1\n",
						},
						{
							Input:  "4\n2 1 4 3\n",
							Output: "1\n",
						},
					},
					Arg: Args{"101189", "c", "gym", ""},
				},
				{
					Name:        "D. Arpa’s letter-marked tree and Mehrdad’s Dokhtar-kosh paths(Hard)",
					TimeLimit:   "2 seconds",
					MemoryLimit: "256 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "4\n1 s\n2 a\n3 s\n",
							Output: "3 1 1 0\n",
						},
						{
							Input:  "5\n1 a\n2 z\n1 a\n4 z\n",
							Output: "4 1 0 1 0\n",
						},
					},
					Arg: Args{"101189", "d", "gym", ""},
				},
			},
			wantErr: false,
		},
		{
			name: "Test #4",
			arg:  Args{"102391", "g", "gym", ""},
			want: []Problem{
				{
					Name:        "G. Lexicographically Minimum Walk",
					TimeLimit:   "2 seconds",
					MemoryLimit: "1024 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "3 3 1 3\n1 2 1\n2 3 7\n1 3 5\n",
							Output: "1 7\n",
						},
						{
							Input:  "3 4 1 3\n1 2 1\n2 1 2\n2 3 7\n1 3 5\n",
							Output: "TOO LONG\n",
						},
						{
							Input:  "2 0 2 1\n",
							Output: "IMPOSSIBLE\n",
						},
					},
					Arg: Args{"102391", "g", "gym", ""},
				},
			},
			wantErr: false,
		},
		{
			name: "Test #5",
			arg:  Args{"283855", "", "group", "bK73bvp3d7"},
			want: []Problem{
				{
					Name:        "A. Buggy Robot",
					TimeLimit:   "2 seconds",
					MemoryLimit: "256 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "4\nLDUR\n",
							Output: "4\n",
						},
						{
							Input:  "5\nRRRUU\n",
							Output: "0\n",
						},
						{
							Input:  "6\nLLRRRR\n",
							Output: "4\n",
						},
					},
					Arg: Args{"283855", "a", "group", "bK73bvp3d7"},
				},
				{
					Name:        "B. Two Cakes",
					TimeLimit:   "1 second",
					MemoryLimit: "256 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "5 2 3\n",
							Output: "1\n",
						},
						{
							Input:  "4 7 10\n",
							Output: "3\n",
						},
					},
					Arg: Args{"283855", "b", "group", "bK73bvp3d7"},
				},
				{
					Name:        "C. Odd sum",
					TimeLimit:   "1 second",
					MemoryLimit: "256 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "4\n-2 2 -3 1\n",
							Output: "3\n",
						},
						{
							Input:  "3\n2 -5 -3\n",
							Output: "-1\n",
						},
					},
					Arg: Args{"283855", "c", "group", "bK73bvp3d7"},
				},
				{
					Name:        "D. Number of Ways",
					TimeLimit:   "2 seconds",
					MemoryLimit: "256 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "5\n1 2 3 0 3\n",
							Output: "2\n",
						},
						{
							Input:  "4\n0 1 -1 0\n",
							Output: "1\n",
						},
						{
							Input:  "2\n4 1\n",
							Output: "0\n",
						},
					},
					Arg: Args{"283855", "d", "group", "bK73bvp3d7"},
				},
				{
					Name:        "E. Propagating tree",
					TimeLimit:   "2 seconds",
					MemoryLimit: "256 megabytes",
					InpStream:   "standard input",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "5 5\n1 2 1 1 2\n1 2\n1 3\n2 4\n2 5\n1 2 3\n1 1 2\n2 1\n2 2\n2 4\n",
							Output: "3\n3\n0\n",
						},
					},
					Arg: Args{"283855", "e", "group", "bK73bvp3d7"},
				},
			},
			wantErr: false,
		},
		{
			name: "Test #6",
			arg:  Args{"277493", "t", "group", "MEqF8b6wBT"},
			want: []Problem{
				{
					Name:        "T. Rhombuses Inside Rectangle",
					TimeLimit:   "2 seconds",
					MemoryLimit: "256 megabytes",
					InpStream:   "rect.in",
					OutStream:   "standard output",
					SampleTests: []SampleTest{
						{
							Input:  "3\n1 1\n2 2\n2 3\n",
							Output: "0\n1\n2\n",
						},
					},
					Arg: Args{"277493", "t", "group", "MEqF8b6wBT"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.GetProblems()
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.GetProblems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Args.GetProblems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgs_SubmitSolution(t *testing.T) {
	source := genRandomString(30)

	type args struct {
		langID string
		source string
	}
	tests := []struct {
		name    string
		arg     Args
		args    args
		wantErr bool
	}{
		{
			name:    "Test #1",
			arg:     Args{"5", "a", "contest", ""},
			args:    args{"54", source},
			wantErr: false,
		},
		{
			name:    "Test #2",
			arg:     Args{"5", "a", "contest", ""},
			args:    args{"54", source},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.arg.SubmitSolution(tt.args.langID, tt.args.source); (err != nil) != tt.wantErr {
				t.Errorf("Args.SubmitSolution() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}