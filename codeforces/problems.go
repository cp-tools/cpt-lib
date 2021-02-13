package codeforces

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type (
	// SampleTest holds sample test case data.
	SampleTest struct {
		Input  string
		Output string
	}

	// Problem holds data of problem.
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

func (p *page) getProblems(arg Args) ([]Problem, error) {
	pd := p.parse()

	problems := make([]Problem, 0)

	problemsTable := pd.Find(`.problemindexholder`)
	problemsTable.Each(func(_ int, row *goquery.Selection) {
		var problem Problem

		problem.Arg, _ = Parse(arg.Group + arg.Contest + row.AttrOr("problemindex", ""))

		// Extract sample test cases of problem.
		for i, sampleInput, sampleOutput := 0, row.Find(`.input>pre`),
			row.Find(`.output>pre`); i < sampleInput.Length(); i++ {

			inpStr := p.MustEval(
				fmt.Sprintf("document.querySelector(\"#%v\").innerText",
					sampleInput.Eq(i).AttrOr(`id`, ``))).String()

			outStr := p.MustEval(
				fmt.Sprintf("document.querySelector(\"#%v\").innerText",
					sampleOutput.Eq(i).AttrOr(`id`, ``))).String()

			problem.SampleTests = append(problem.SampleTests, SampleTest{
				Input: inpStr, Output: outStr,
			})
		}

		header := row.Find(`.header`)
		// Bulk extract rest of the data from header.
		problem.Name = clean(header.Find(`.title`).Text())
		problem.TimeLimit = clean(header.Find(`.time-limit`).Contents().Last().Text())
		problem.MemoryLimit = clean(header.Find(`.memory-limit`).Contents().Last().Text())
		problem.InpStream = clean(header.Find(".input-file").Contents().Last().Text())
		problem.OutStream = clean(header.Find(".output-file").Contents().Last().Text())

		problems = append(problems, problem)
	})

	return problems, nil
}

// GetProblems returns problem(s) meta data, along with sample tests.
//
// If the problem is not specified, returns data of all problems in
// the specified contest. In some older contests, the complete
// problemset page is not present. In such cases, fallback to running
// GetProblems() for each problem in the contest.
//
// SolveStatus and SolveCount are not parsed by this.
// Use GetDashboard() if you require these fields.
func (arg Args) GetProblems() ([]Problem, error) {
	link, err := arg.ProblemsPage()
	if err != nil {
		return nil, err
	}

	p, err := loadPage(link)
	if err != nil {
		return nil, err
	}
	defer p.Close()

	if _, err := p.Race().Element(`#jGrowl .message`).Handle(handleErrMsg).
		Element(`.problemindexholder`).Do(); err != nil {
		return nil, err
	}

	// Wait till all problems have loaded.
	p.WaitLoad()

	return p.getProblems(arg)
}

// SubmitSolution submits given file to the judging server,
// and returns a channel on a successful submission.
// The channel contains the live status of the submission.
// View GetSubmissions() for more details on the returned channel.
//
// langName is the codeforces configured language to use. See the
// variable map LanguageID for the list of supported languages.
func (arg Args) SubmitSolution(langName string, file string) (<-chan Submission, error) {
	// problem not specified, return invalid
	if arg.Problem == "" {
		return nil, ErrInvalidSpecifier
	}

	if _, ok := LanguageID[langName]; !ok {
		return nil, fmt.Errorf("invalid language")
	}

	// check if given file exists
	if fl, err := os.Stat(file); os.IsNotExist(err) || fl.IsDir() {
		return nil, fmt.Errorf("invalid file path")
	}

	link, err := arg.ProblemsPage()
	if err != nil {
		return nil, err
	}

	p, err := loadPage(link)
	if err != nil {
		return nil, err
	}

	if _, err := p.Race().Element(`#jGrowl .message`).Handle(handleErrMsg).
		Element(`#footer`).Do(); err != nil {
		p.Close()
		return nil, err
	}

	// Check if user is logged in.
	if !p.MustHas(`#header a[href^="/profile/"]`) {
		p.Close()
		return nil, fmt.Errorf("no logged in session present")
	}

	// Check if submitting is possible at all.
	if !p.MustHas(`input.submit`) {
		p.Close()
		return nil, fmt.Errorf("problem not open for submission")
	}

	// Check if specified language can be selected.
	// If this is allowed, so is submitting.
	if !p.MustHasR(`select>option[value]`, regexp.QuoteMeta(langName)) {
		p.Close()
		return nil, fmt.Errorf("language not allowed in problem")
	}

	// All cases have been handled. Submit the solution.
	p.MustElement(`select[name="programTypeId"]`).MustSelect(langName)
	p.MustElement(`input[name="sourceFile"]`).MustSetFiles(file)
	p.MustElement(`input.submit`).MustClick().WaitInvisible()

	if _, err := p.Race().Element(`.error`).Handle(handleErrMsg).
		Element(`tr[data-submission-id]`).Do(); err != nil {
		// Example error message: "exact submission done before"
		p.Close()
		return nil, err
	}

	// Realtime verdict of submission.
	chanSubmission := make(chan Submission)
	go func() {
		defer p.Close()
		defer close(chanSubmission)

		for timer := time.Now(); ; time.Sleep(time.Millisecond * 400) {
			submissions, _ := p.getSubmissions(arg)
			if len(submissions) == 0 {
				break
			}

			chanSubmission <- submissions[0]
			if !submissions[0].IsJudging {
				break
			}

			if time.Since(timer) > 2*time.Second {
				// Reload the page every 2 seconds.
				// This is to handle websocket failure
				// and completion of judgement in WA case.
				p.MustReload().MustWaitLoad()
				timer = time.Now()
			}
		}
	}()

	return chanSubmission, nil
}
