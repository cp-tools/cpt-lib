package codeforces

import (
	"fmt"
	"os"
	"regexp"

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
// file is the submissions file to upload on the form.
//
// If submission completed successfully, returns nil error.
func (arg Args) SubmitSolution(langName string, file string) error {
	// problem not specifed, return invalid
	if len(arg.Contest) == 0 || len(arg.Problem) == 0 {
		return ErrInvalidSpecifier
	}

	if _, ok := LanguageID[langName]; !ok {
		return fmt.Errorf("Invalid language")
	}

	// check if given file exists
	if fl, err := os.Stat(file); os.IsNotExist(err) || fl.IsDir() {
		return fmt.Errorf("Invalid file path")
	}

	link := arg.ProblemsPage()
	page, msg, err := loadPage(link)
	if err != nil {
		return err
	}
	defer page.Close()

	if msg != "" {
		return fmt.Errorf(msg)
	}

	// check if user is logged in
	if !page.MustHas(selCSSHandle) {
		return fmt.Errorf("No logged in session present")
	}

	// check if specified language can be selected
	if !page.MustHasMatches(`select>option[value]`, regexp.QuoteMeta(langName)) {
		return fmt.Errorf("Language not supported")
	}

	// do the submitting here! (really simple)
	page.MustElement(`select[name="programTypeId"]`).MustSelect(langName)
	page.MustElement(`input[name="sourceFile"]`).MustSetFiles(file)
	page.MustElement(`input.submit`).MustClick()

	elm := page.MustElement(selCSSError, `tr[data-submission-id]`)

	if elm.MustMatches(selCSSError) {
		// static error message (exact submission done before)
		msg := clean(elm.MustText())
		return fmt.Errorf(msg)
	}
	return nil
}
