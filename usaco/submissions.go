package usaco

import (
	"regexp"
)

type (
	// TestCaseVerdict holds information about
	// verdict of an individual test case.
	TestCaseVerdict struct {
		Index   int
		Verdict string
		Time    string
		Memory  string
	}

	// Verdict holds verdict of last submission.
	Verdict struct {
		// compilation error, file error is parsed here.
		// Is empty string if no errors occurred.
		Status      string
		LastVerdict chan TestCaseVerdict
	}
)

// Verdict Status types
const (
	VerdictStatusOK = "OK" // submission compiled successfully
	VerdictStatusNA = "NA" // no submission found
	VerdictStatusCE = "CE" // compilation error
	VerdictStatusFS = "FS" // Failed sample test case
)

// GetSubmission parses and returns submission status (verdict)
// of latest submission in problem.
//
// A buffered channel is opened until all tests have been judged.
// Read verdicts of latest testcase from 'Verdict.LastVerdict'.
func (arg Args) GetSubmission() (Verdict, error) {
	if arg.Cpid == "" {
		return Verdict{}, ErrInvalidSpecifier
	}

	link := arg.ProblemPage()
	page, err := loadPage(link)
	if err != nil {
		return Verdict{}, err
	}
	page.MustElement(`#last-status[data-sid]`)
	doc := processHTML(page)
	// page close after goroutine completes
	var verdict Verdict

	lastStatusElm := doc.Find("#last-status")
	if lastStatusElm.AttrOr("data-sid", "-1") == "-1" {
		verdict.Status = VerdictStatusNA
		return verdict, nil
	}

	if lastStatusElm.AttrOr("class", "") == "status-no" {
		// check if compilation error (or failed sample)
		str := lastStatusElm.Text()
		if match, _ := regexp.MatchString(`Compilation Error`, str); match {
			verdict.Status = VerdictStatusCE

		} else {
			verdict.Status = VerdictStatusFS
		}
		return verdict, nil
	}

	verdict.Status = VerdictStatusOK
	// create buffered channel for latest testcase verdict
	verdict.LastVerdict = make(chan TestCaseVerdict, 100)
	go func() {
		defer page.Close()
		defer close(verdict.LastVerdict)

		// parse 'index' test case verdict
		index := 1

		for true {
			doc := processHTML(page)
			// order of statements below matters
			judgeDone := page.MustHasMatches(`#last-status>p`,
				`Results below show the outcome`)
			testcases := doc.Find(`#trial-information .masterTooltip`)

			for index <= testcases.Length() {
				// found 'index' judgement
				testcase := testcases.Eq(index - 1)
				tcVerdict := TestCaseVerdict{
					Index:   index,
					Verdict: testcase.AttrOr("title", ""),
					Memory:  testcase.Find(".info>span").Eq(0).Text(),
					Time:    testcase.Find(".info>span").Eq(1).Text(),
				}
				verdict.LastVerdict <- tcVerdict
				index++
			}

			if judgeDone == true {
				// all data parsed; break
				break
			}
		}
	}()
	return verdict, nil
}
