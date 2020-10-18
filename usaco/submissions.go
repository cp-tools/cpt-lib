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
func (arg Args) GetSubmission() (<-chan TestCaseVerdict, string, error) {
	if arg.Cpid == "" {
		return nil, "", ErrInvalidSpecifier
	}

	link := arg.ProblemPage()
	page, err := loadPage(link)
	if err != nil {
		return nil, "", err
	}
	page.MustElement(`#last-status[data-sid]`)
	doc := processHTML(page)

	lastStatusElm := doc.Find("#last-status")
	if lastStatusElm.AttrOr("data-sid", "-1") == "-1" {
		return nil, VerdictStatusNA, nil
	}

	if lastStatusElm.AttrOr("class", "") == "status-no" {
		// check if compilation error (or failed sample)
		str := lastStatusElm.Text()
		if match, _ := regexp.MatchString(`Compilation Error`, str); match {
			return nil, VerdictStatusCE, nil
		}
		return nil, VerdictStatusFS, nil
	}

	// create buffered channel for latest testcase verdict
	chanTestCaseVerdict := make(chan TestCaseVerdict, 100)
	go func() {
		defer page.Close()
		defer close(chanTestCaseVerdict)

		// parse 'index' test case verdict
		index := 1

		for true {
			doc := processHTML(page)
			// order of statements below matters
			judgeDone := page.MustHasR(`#last-status>p`,
				`Results below show the outcome`)
			testcases := doc.Find(`#trial-information .masterTooltip`)

			for index <= testcases.Length() {
				// found 'index' judgement
				testcase := testcases.Eq(index - 1)
				testCaseVerdict := TestCaseVerdict{
					Index:   index,
					Verdict: testcase.AttrOr("title", ""),
					Memory:  testcase.Find(".info>span").Eq(0).Text(),
					Time:    testcase.Find(".info>span").Eq(1).Text(),
				}
				chanTestCaseVerdict <- testCaseVerdict
				index++
			}

			if judgeDone == true {
				// all data parsed; break
				break
			}
		}
	}()
	return chanTestCaseVerdict, VerdictStatusOK, nil
}
