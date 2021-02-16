package codeforces

import (
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type (
	// Submission holds submission data.
	Submission struct {
		ID            string
		When          time.Time
		Who           string
		Problem       string
		Language      string
		Verdict       string
		VerdictStatus int
		Time          string
		Memory        string
		IsJudging     bool
		Arg           Args
	}
)

// Submissions verdict status.
const (
	VerdictAC = 1 // Accepted

	VerdictWA  = 2 // Wrong Answer
	VerdictRTE = 4 // Run Time Error

	VerdictCE  = 5 // Compilation Error
	VerdictTLE = 6 // Time Limit Exceeded
	VerdictMLE = 7 // Memory Limit Exceeded
	VerdictILE = 8 // Idleness Limit Exceeded

	VerdictDOJ  = 9  // Denial Of Judgement
	VerdictSkip = 10 // Skipped
	VerdictHack = 11 // Hacked

	// Depreciated. Use VerdictAC instead.
	VerdictPretestPass = 12 // Pretests passed
)

func (p *page) getSubmissions(arg Args) ([]Submission, error) {
	pd := p.parse()

	submissions := make([]Submission, 0)

	submissionTableRows := pd.Find(`tr[data-submission-id]`)
	submissionTableRows.Each(func(_ int, row *goquery.Selection) {
		var submission Submission

		submission.Arg, _ = Parse(hostURL + row.Find(`td`).Eq(3).Find(`a`).AttrOr(`href`, ``))
		if arg.Problem != "" && arg.Problem != submission.Arg.Problem {
			return
		}

		row.Find(`td`).Each(func(cellIndex int, cell *goquery.Selection) {
			switch cellIndex {
			case 0:
				submission.ID = clean(cell.Text())

			case 1:
				submission.When = parseTime(cell.Text())

			case 2:
				submission.Who = clean(cell.Text())

			case 3:
				submission.Problem = clean(cell.Text())

			case 4:
				submission.Language = clean(cell.Text())

			case 5:
				submission.Verdict = clean(cell.Text())

				verdictMap := map[string]int{
					"OK":                      VerdictAC,
					"WRONG_ANSWER":            VerdictWA,
					"RUNTIME_ERROR":           VerdictRTE,
					"COMPILATION_ERROR":       VerdictCE,
					"TIME_LIMIT_EXCEEDED":     VerdictTLE,
					"MEMORY_LIMIT_EXCEEDED":   VerdictMLE,
					"IDLENESS_LIMIT_EXCEEDED": VerdictILE,
					"CRASHED":                 VerdictDOJ,
					"SKIPPED":                 VerdictSkip,
					"CHALLENGED":              VerdictHack,
				}

				verdictStatus := cell.Find(`.submissionVerdictWrapper`).
					AttrOr(`submissionverdict`, ``)
				if v, ok := verdictMap[verdictStatus]; ok {
					submission.VerdictStatus = v
					submission.IsJudging = false
				} else if strings.Contains(submission.Verdict, "Compilation error") {
					// Compilation error verdict is updated only on page reload.
					// Resorted to manually handling this bug on codeforces.
					submission.VerdictStatus = VerdictCE
					submission.IsJudging = false
				} else {
					submission.IsJudging = true
				}

			case 6:
				submission.Time = clean(cell.Text())

			case 7:
				submission.Memory = clean(cell.Text())
			}
		})

		submissions = append(submissions, submission)
	})

	return submissions, nil
}

// GetSubmissions returns submissions metadata of given user.
// If contest is not specified, returns all submissions of user.
//
// Due to a bug on codeforces, fetching submissions in group contests
// are not supported, when the contest isn't specified.
//
// Set pageCount to maximum number of pages to parse. Each page consists of 50
// rows of data. If pageCount is 1, the returned channel will keep returning page
// data, till all verdicts of submissions in the page are declared.
func (arg Args) GetSubmissions(handle string, pageCount uint) (<-chan []Submission, error) {
	link, err := arg.SubmissionsPage(handle)
	if err != nil {
		return nil, err
	}

	p, err := loadPage(link)
	if err != nil {
		return nil, err
	}

	if _, err := p.Race().Element(`#jGrowl .message`).Handle(handleErrMsg).
		Element(`tr[data-submission-id]`).Do(); err != nil {
		p.Close()
		return nil, err
	}

	if p.MustInfo().URL != link {
		p.Close()
		// An unexpected redirect occurred.
		// Return error notification.
		return nil, handleErrMsg(p.MustElement(`#jGrowl .message`))
	}

	// Wait till alls rows are loaded.
	p.MustWaitLoad()

	// @todo Add support for excluding unofficial submissions

	// create buffered channel for submissions
	chanSubmissions := make(chan []Submission)
	go func() {
		defer p.Close()
		defer close(chanSubmissions)

		// Only one page to parse. Keep parsing till all verdicts are declared.
		if pageCount == 1 {
			for true {
				// Keep parsing verdict till
				// all submission verdicts are finalised.
				submissions, _ := p.getSubmissions(arg)
				chanSubmissions <- submissions

				IsJudging := false
				for _, sub := range submissions {
					IsJudging = (IsJudging || sub.IsJudging)
				}
				if !IsJudging {
					break
				}

				time.Sleep(time.Millisecond * 400)
			}
		} else {
			// Parse each page (without waiting for judgement to complete).
			for ; pageCount > 0; pageCount-- {
				// Ignore error, write whatever is parsed.
				submissions, _ := p.getSubmissions(arg)
				chanSubmissions <- submissions

				if !p.MustHasR(`.pagination li>a`, `→`) || pageCount == 1 {
					// All pages parsed.
					break
				}

				// Move to the next page (click the next button).
				p.MustElementR(`.pagination li>a`, `→`).MustClick().WaitInvisible()
				p.WaitLoad()
			}
		}
	}()
	return chanSubmissions, nil
}

// GetSourceCode returns submission code of given submission.
func (sub Submission) GetSourceCode() (string, error) {
	link, err := sub.SourceCodePage()
	if err != nil {
		return "", err
	}

	p, err := loadPage(link)
	if err != nil {
		return "", err
	}
	defer p.Close()

	if _, err := p.Race().Element(`#jGrowl .message`).Handle(handleErrMsg).
		Element(`#program-source-text`).Do(); err != nil {
		return "", err
	}

	if p.MustInfo().URL != link {
		// An unexpected redirect occurred.
		// Return error notification.
		return "", handleErrMsg(p.MustElement(`#jGrowl .message`))
	}

	sourceCode := p.MustEval(`Codeforces.filterClipboardText(
		document.querySelector("#program-source-text").innerText)`).String()
	return sourceCode, nil
}
