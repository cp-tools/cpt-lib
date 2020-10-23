package codeforces

import (
	"fmt"
	"os"
	"regexp"
	"time"

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
func (arg Args) ProblemsPage() (link string, err error) {
	if arg.Contest == "" {
		return "", ErrInvalidSpecifier
	}

	switch arg.Class {
	case ClassGroup:
		if arg.Group == "" {
			return "", ErrInvalidSpecifier
		}

		if arg.Problem == "" {
			link = fmt.Sprintf("%v/group/%v/contest/%v/problems", hostURL, arg.Group, arg.Contest)
		} else {
			link = fmt.Sprintf("%v/group/%v/contest/%v/problem/%v", hostURL, arg.Group, arg.Contest, arg.Problem)
		}

	case ClassContest, ClassGym:
		if arg.Problem == "" {
			link = fmt.Sprintf("%v/%v/%v/problems", hostURL, arg.Class, arg.Contest)
		} else {
			link = fmt.Sprintf("%v/%v/%v/problem/%v", hostURL, arg.Class, arg.Contest, arg.Problem)
		}

	default:
		return "", ErrInvalidSpecifier
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

	link, err := arg.ProblemsPage()
	if err != nil {
		return nil, err
	}

	page, msg, err := loadPage(link)
	if err != nil {
		return nil, err
	}
	defer page.Close()

	if msg != "" {
		return nil, fmt.Errorf(msg)
	}

	doc := processHTML(page)

	// to hold problem data
	var probs []Problem
	table := doc.Find(".problemindexholder")
	table.Each(func(_ int, row *goquery.Selection) {
		probArg, _ := Parse(arg.Group + arg.Contest + row.AttrOr("problemindex", ""))

		// sample tests of problem
		var sampleTests []SampleTest
		inp, out := row.Find(".input"), row.Find(".output")
		for i := 0; i < inp.Length(); i++ {
			inpStr, _ := inp.Find("pre").Eq(i).Html()
			outStr, _ := out.Find("pre").Eq(i).Html()
			sampleTests = append(sampleTests, SampleTest{
				Input:  clean(inpStr) + "\n",
				Output: clean(outStr) + "\n",
			})
		}

		header := row.Find(".header")
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
// langName is codeforces specified name of language to submit in.
// file is the submission code file to upload.
//
// If submitted successfully, returns a chan updating verdict
// of the current submission. View GetSubmissions() for more details.
func (arg Args) SubmitSolution(langName string, file string) (<-chan Submission, error) {
	// problem not specifed, return invalid
	if arg.Problem == "" {
		return nil, ErrInvalidSpecifier
	}

	if _, ok := LanguageID[langName]; !ok {
		return nil, fmt.Errorf("Invalid language")
	}

	// check if given file exists
	if fl, err := os.Stat(file); os.IsNotExist(err) || fl.IsDir() {
		return nil, fmt.Errorf("Invalid file path")
	}

	link, err := arg.ProblemsPage()
	if err != nil {
		return nil, err
	}

	page, msg, err := loadPage(link)
	if err != nil {
		return nil, err
	}

	if msg != "" {
		defer page.Close()
		return nil, fmt.Errorf(msg)
	}

	// check if user is logged in
	if !page.MustHas(selCSSHandle) {
		defer page.Close()
		return nil, fmt.Errorf("No logged in session present")
	}

	// check if submitting is possible at all.
	if !page.MustHas(`input.submit`) {
		defer page.Close()
		return nil, fmt.Errorf("Problem not open for submission")
	}

	// check if specified language can be selected
	// if this is allowed, so is submitting.
	if !page.MustHasR(`select>option[value]`, regexp.QuoteMeta(langName)) {
		defer page.Close()
		return nil, fmt.Errorf("Language not allowed in problem")
	}

	// do the submitting here! (really simple)
	page.MustElement(`select[name="programTypeId"]`).MustSelect(langName)
	page.MustElement(`input[name="sourceFile"]`).MustSetFiles(file)
	page.MustElement(`input.submit`).MustClick().WaitInvisible()
	page.MustWaitLoad()

	if page.MustHas(selCSSError) {
		defer page.Close()
		// static error message (exact submission done before)
		elm := page.MustElement(selCSSError)
		msg := clean(elm.MustText())
		return nil, fmt.Errorf(msg)

	}

	// return live progress of submission
	chanSubmission := make(chan Submission, 500)
	go func() {
		defer page.Close()
		defer close(chanSubmission)

		for true {
			submissions, _ := arg.parseSubmissions(page)
			chanSubmission <- submissions[0]
			if submissions[0].IsJudging == false {
				break
			}
			time.Sleep(time.Millisecond * 350)
		}
	}()

	return chanSubmission, nil
}
