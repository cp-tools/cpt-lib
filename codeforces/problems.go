package codeforces

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type (
	// SampleTest maps sample input to sample output.
	SampleTest struct {
		Input  string
		Output string
	}

	// Problem data is parsed to this struct.
	Problem struct {
		Name        string
		TimeLimit   string
		MemoryLimit string
		InpStream   string
		OutStream   string
		SampleTests []SampleTest
		SolveCount  int
		SolveStatus int
		Arg         Args
	}
)

// Different values of 'SolveStatus'.
const (
	SolveAccepted     = 1
	SolveRejected     = 0
	SolveNotAttempted = -1
)

// ProblemsPage returns link to problem(s) page in contest
func (arg Args) ProblemsPage() (link string) {
	// problem specified
	if len(arg.Problem) != 0 {
		if arg.Class == ClassGroup {
			link = fmt.Sprintf("%v/group/%v/contest/%v/problem/%v",
				hostURL, arg.Group, arg.Contest, arg.Problem)
		} else {
			link = fmt.Sprintf("%v/%v/%v/problem/%v",
				hostURL, arg.Class, arg.Contest, arg.Problem)
		}
	} else {
		if arg.Class == ClassGroup {
			link = fmt.Sprintf("%v/group/%v/contest/%v/problems",
				hostURL, arg.Group, arg.Contest)
		} else {
			link = fmt.Sprintf("%v/%v/%v/problems",
				hostURL, arg.Class, arg.Contest)
		}
	}
	return
}

// GetProblems parses problem(s) details along with sample tests.
// If problem field is not specified, extracts details of all problems
// in the contest.
//
// In some older contests, complete problemset page is not supported.
// Preferably fallback to parsing individual problems if entire parsing
// fails.
//
// Doesn't fetch 'SolveStatus' and 'SolveCount' of problem.
// Use GetDashboard() to fetch these info fields.
func (arg Args) GetProblems() ([]Problem, error) {
	if len(arg.Contest) == 0 {
		return nil, ErrInvalidSpecifier
	}

	link := arg.ProblemsPage()
	resp, err := SessCln.Get(link)
	if err != nil {
		return nil, err
	}
	body, msg := parseResp(resp)
	if len(msg) != 0 {
		// shouldn't return any error if success
		return nil, fmt.Errorf(msg)
	}
	// to hold problem data
	var probs []Problem

	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	table := doc.Find(".problemindexholder")
	table.Each(func(_ int, prob *goquery.Selection) {
		probArg, _ := Parse(arg.Group + arg.Contest + prob.AttrOr("problemindex", ""))

		// sample tests of problem
		var sampleTests []SampleTest
		inp, out := prob.Find(".input"), prob.Find(".output")
		for i := 0; i < inp.Length(); i++ {
			inpStr, _ := inp.Find("pre").Eq(i).Html()
			outStr, _ := out.Find("pre").Eq(i).Html()
			sampleTests = append(sampleTests, SampleTest{
				Input:  clean(inpStr) + "\n",
				Output: clean(outStr) + "\n",
			})
		}

		header := prob.Find(".header")
		probs = append(probs, Problem{
			Name:        getText(header, ".title"),
			TimeLimit:   clean(header.Find(".time-limit").Contents().Last().Text()),
			MemoryLimit: clean(header.Find(".memory-limit").Contents().Last().Text()),
			InpStream:   clean(header.Find(".input-file").Contents().Last().Text()),
			OutStream:   clean(header.Find(".output-file").Contents().Last().Text()),
			SampleTests: sampleTests,
			Arg:         probArg,
		})
	})
	return probs, nil
}

// SubmitSolution submits source code to specifed problem.
// langID is codeforces specified id of language to submit in.
// View cp-tools/codeforces.wiki for list of valid ID's.
// Source is code text to submit.
//
// If submission completed successfully, returns nil error.
func (arg Args) SubmitSolution(langID string, source string) error {
	// problem not specifed, return invalid
	if len(arg.Contest) == 0 || len(arg.Problem) == 0 {
		return ErrInvalidSpecifier
	}
	// if langID invalid, return invalid
	isLangIDValid := false
	for _, v := range LanguageID {
		if v == langID {
			isLangIDValid = true
			break
		}
	}
	if isLangIDValid == false {
		return fmt.Errorf("Invalid language id")
	}

	link := arg.ProblemsPage()
	resp, err := SessCln.Get(link)
	if err != nil {
		return err
	}
	body, msg := parseResp(resp)
	if len(msg) != 0 {
		// shouldn't return any error if success
		return fmt.Errorf(msg)
	}

	// hidden form data
	csrf := findCsrf(body)
	ftaa := genRandomString(18)
	bfaa := genRandomString(32)

	resp, err = SessCln.PostForm(link, url.Values{
		"csrf_token":            {csrf},
		"ftaa":                  {ftaa},
		"bfaa":                  {bfaa},
		"action":                {"submitSolutionFormSubmitted"},
		"submittedProblemIndex": {arg.Problem},
		"programTypeId":         {langID},
		"contestId":             {arg.Contest},
		"source":                {source},
		"tabSize":               {"4"},
		"_tta":                  {"176"},
		"sourceCodeConfirmed":   {"true"},
	})
	if err != nil {
		return err
	}

	body, msg = parseResp(resp)
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if msg1 := getText(doc.Selection, ".error.for__source"); len(msg1) != 0 {
		// static error message (exact submission done before)
		return fmt.Errorf(msg1)
	}
	// successful submission should have message :
	// "Solution to the problem X has been submitted successfully"
	if strings.EqualFold(msg, "Solution to the problem "+arg.Problem+" has been submitted successfully") {
		// submission successful!!!
		return nil
	}

	return fmt.Errorf(msg)
}
