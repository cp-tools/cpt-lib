package usaco

import (
	"reflect"
	"testing"
)

func TestArgs_GetSubmission(t *testing.T) {
	type fields struct {
		Cpid string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []TestCaseVerdict
		want1   string
		wantErr bool
	}{
		{
			name:   "Test #1",
			fields: fields{"545"},
			want: []TestCaseVerdict{
				{
					Index:   1,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "9ms",
				},
				{
					Index:   2,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "1ms",
				},
				{
					Index:   3,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "1ms",
				},
				{
					Index:   4,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "2ms",
				},
				{
					Index:   5,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "3ms",
				},
				{
					Index:   6,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "1ms",
				},
				{
					Index:   7,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "9ms",
				},
				{
					Index:   8,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "19ms",
				},
				{
					Index:   9,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "27ms",
				},
				{
					Index:   10,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "28ms",
				},
				{
					Index:   11,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "23ms",
				},
				{
					Index:   12,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "19ms",
				},
				{
					Index:   13,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "27ms",
				},
				{
					Index:   14,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "19ms",
				},
				{
					Index:   15,
					Verdict: "Correct answer",
					Memory:  "1.2mb",
					Time:    "23ms",
				},
			},
			want1:   "OK",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg := Args{
				Cpid: tt.fields.Cpid,
			}
			got, got1, err := arg.GetSubmission()
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.GetSubmission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got1 != tt.want1 {
				t.Errorf("Args.GetSubmission() got1 = %v, want %v", got1, tt.want1)
			}

			val := make([]TestCaseVerdict, 0)
			for v := range got {
				val = append(val, v)
			}

			if !reflect.DeepEqual(val, tt.want) {
				t.Errorf("Args.GetSubmission() got = %v, want %v", val, tt.want)
			}
		})
	}
}
