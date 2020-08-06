package codeforces

import (
	"fmt"
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
	page := Browser.Page(link).WaitLoad()

	if msg := cE(page); msg != "" {
		// shouldn't return any error if success
		return nil, fmt.Errorf(msg)
	}
	body := page.Element("html").HTML()

	// to hold problem data
	var probs []Problem

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(body))
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
// file is the submissions file to upload on the form.
//
// If submission completed successfully, returns nil error.
func (arg Args) SubmitSolution(langID string, file string) error {
	// problem not specifed, return invalid
	if len(arg.Contest) == 0 || len(arg.Problem) == 0 {
		return ErrInvalidSpecifier
	}
	// if langID invalid, return invalid
	langIDName := ""
	for langName, v := range LanguageID {
		if v == langID {
			langIDName = langName
			break
		}
	}
	if langIDName == "" {
		return fmt.Errorf("Invalid language id")
	}

	link := arg.ProblemsPage()
	page := Browser.Page(link).WaitLoad()

	if msg := cE(page); msg != "" {
		// shouldn't return any error if success
		return fmt.Errorf(msg)
	}

	// do the submitting here! (really simple)
	page.Element(`select[name="programTypeId"]`).Select(langIDName)
	page.Element(`input[name="sourceFile"]`).SetFiles(file)
	page.Element(`input.submit`).Click().WaitLoad()

	if msg := page.Elements(`.error.for__source`); len(msg) != 0 {
		// static error message (exact submission done before)
		return fmt.Errorf(msg.First().Text())
	}

	msg := cE(page)
	if strings.EqualFold(msg, "Solution to the problem "+arg.Problem+" has been submitted successfully") {
		// submission successful!!
		return nil
	}

	return fmt.Errorf(msg)
}
