package usaco

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestArgs_ProblemPage(t *testing.T) {
	type fields struct {
		Cpid string
	}
	tests := []struct {
		name     string
		fields   fields
		wantLink string
	}{
		{
			name: "Test #1",
			fields: fields{
				Cpid: "246",
			},
			wantLink: "http://usaco.org/index.php?page=viewproblem2&cpid=246",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg := Args{
				Cpid: tt.fields.Cpid,
			}
			if gotLink := arg.ProblemPage(); gotLink != tt.wantLink {
				t.Errorf("Args.ProblemPage() = %v, want %v", gotLink, tt.wantLink)
			}
		})
	}
}

func TestArgs_GetProblem(t *testing.T) {
	type fields struct {
		Cpid string
	}
	tests := []struct {
		name    string
		fields  fields
		want    Problem
		wantErr bool
	}{
		{
			name:   "Test #1",
			fields: fields{"246"},
			want: Problem{
				Name:      "Milk Scheduling",
				Contest:   "USACO 2013 February Contest, Silver",
				InpStream: `msched.in`,
				OutStream: `msched.out`,
				SampleTests: []SampleTest{
					{
						Input:  "3 1\n10\n5\n6\n3 2\n",
						Output: "11\n",
					},
				},
				Arg: Args{"246"},
			},
			wantErr: false,
		},
		{
			name:   "Test #2",
			fields: fields{"994"},
			want: Problem{
				Name:      "Farmer John Solves 3SUM",
				Contest:   "USACO 2020 January Contest, Gold",
				InpStream: "threesum.in",
				OutStream: "threesum.out",
				SampleTests: []SampleTest{
					{
						Input:  "7 3\n2 0 -1 1 -2 3 3\n1 5\n2 4\n1 7\n",
						Output: "2\n1\n4\n",
					},
				},
				Arg: Args{"994"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg := Args{
				Cpid: tt.fields.Cpid,
			}
			got, err := arg.GetProblem()
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.GetProblem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Args.GetProblem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgs_SubmitSolution(t *testing.T) {
	sFile, _ := ioutil.TempFile(os.TempDir(), "cpt-submission")
	defer os.Remove(sFile.Name())
	sFile.WriteString(`int main(){return 0;}`)

	type fields struct {
		Cpid string
	}
	type args struct {
		langName string
		file     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Test #1",
			fields:  fields{"941"},
			args:    args{"C++", sFile.Name()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg := Args{
				Cpid: tt.fields.Cpid,
			}
			if err := arg.SubmitSolution(tt.args.langName, tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("Args.SubmitSolution() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
