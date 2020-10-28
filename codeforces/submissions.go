package codeforces

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
)

type (
	// Submission holds submission data.
	Submission struct {
		ID        string
		When      time.Time
		Who       string
		Problem   string
		Language  string
		Verdict   string
		Time      string
		Memory    string
		IsJudging bool
		Arg       Args
	}
)

// SubmissionsPage returns link to user submissions page.
func (arg Args) SubmissionsPage(handle string) (link string, err error) {
	// Contest not specified.
	if arg.Contest == "" {
		if handle == "" {
			// Extract handle from homepage.
			// No actual login is done in below code.
			var err error
			if handle, err = login("", ""); err != nil {
				return "", ErrInvalidSpecifier
			}
		}

		link = fmt.Sprintf("%v/submissions/%v", hostURL, handle)
		return
	}

	switch arg.Class {
	case ClassGroup:
		if handle != "" {
			// Fetching others submissions not possible.
			return "", ErrInvalidSpecifier
		}

		link = fmt.Sprintf("%v/group/%v/contest/%v/my", hostURL, arg.Group, arg.Contest)

	case ClassContest, ClassGym:
		if handle == "" {
			link = fmt.Sprintf("%v/%v/%v/my", hostURL, arg.Class, arg.Contest)
		} else {
			link = fmt.Sprintf("%v/submissions/%v/%v/%v", hostURL, handle, arg.Class, arg.Contest)
		}

	default:
		return "", ErrInvalidSpecifier
	}

	return
}

// SourceCodePage returns link to solution submission code.
func (sub Submission) SourceCodePage() (link string, err error) {
	if sub.ID == "" || sub.Arg.Contest == "" {
		return "", ErrInvalidSpecifier
	}

	switch sub.Arg.Class {
	case ClassGroup:
		link = fmt.Sprintf("%v/group/%v/contest/%v/submission/%v", hostURL, sub.Arg.Group, sub.Arg.Contest, sub.ID)

	case ClassContest, ClassGym:
		link = fmt.Sprintf("%v/%v/%v/submission/%v", hostURL, sub.Arg.Class, sub.Arg.Contest, sub.ID)

	default:
		return "", ErrInvalidSpecifier
	}

	return
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

	page, msg, err := loadPage(link, `tr[data-submission-id]`)
	if err != nil {
		return nil, err
	}

	if msg != "" {
		defer page.Close()
		return nil, fmt.Errorf(msg)
	}

	// @todo Add support for excluding unofficial submissions

	// create buffered channel for submissions
	chanSubmissions := make(chan []Submission, 10)
	go func() {
		defer page.Close()
		defer close(chanSubmissions)

		if pageCount == 1 {
			// loop till 'isDone' is true
			for true {
				submissions, isDone := arg.parseSubmissions(page)
				chanSubmissions <- submissions
				if isDone == true {
					break
				}
				time.Sleep(time.Millisecond * 350)
			}
		} else {
			// iterate till no more valid required pages left
			for ; pageCount > 0; pageCount-- {
				page.WaitLoad()

				submissions, _ := arg.parseSubmissions(page)
				chanSubmissions <- submissions

				if !page.MustHasR(".pagination li a", "→") || pageCount < 2 {
					// no more pages to parse
					break
				}
				// click navigation button and wait till loads
				elm := page.MustElementR(".pagination li a", "→")
				elm.MustClick().WaitInvisible()
				page.MustElement(`tr[data-submission-id]`)
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

	page, msg, err := loadPage(link, `#program-source-text`)
	if err != nil {
		return "", err
	}
	defer page.Close()

	if msg != "" {
		return "", fmt.Errorf(msg)
	}

	// extract source code from html body
	source := page.MustEval(`Codeforces.filterClipboardText(
		document.querySelector("#program-source-text").innerText)`).String()
	return source, nil
}

// parse specified submissions from current page.
func (arg Args) parseSubmissions(page *rod.Page) ([]Submission, bool) {
	submissions := make([]Submission, 0)
	isDone := true

	doc := processHTML(page)
	// WARNING! ugly extraction code ahead. Don't peep XD
	table := doc.Find("tr[data-submission-id]")
	table.Each(func(_ int, row *goquery.Selection) {
		var submissionRow Submission
		// extract contest args from html attr label
		subArg, _ := Parse(hostURL + getAttr(row, "td:nth-of-type(4) a", "href"))
		if arg.Problem != "" && arg.Problem != subArg.Problem {
			return
		}
		submissionRow.Arg = subArg

		row.Find("td").Each(func(cellIdx int, cell *goquery.Selection) {
			switch cellIdx {
			case 0:
				id := clean(cell.Text())
				submissionRow.ID = id

			case 1:
				when := parseTime(clean(cell.Text()))
				submissionRow.When = when

			case 2:
				who := clean(cell.Text())
				submissionRow.Who = who

			case 3:
				problem := clean(cell.Text())
				submissionRow.Problem = problem

			case 4:
				language := clean(cell.Text())
				submissionRow.Language = language

			case 5:
				isJudging := cell.AttrOr("waiting", "") == "true"
				submissionRow.IsJudging = isJudging
				if isJudging == true {
					isDone = false
				}

				verdict := clean(cell.Text())
				submissionRow.Verdict = verdict

			case 6:
				time := clean(cell.Text())
				submissionRow.Time = time

			case 7:
				memory := clean(cell.Text())
				submissionRow.Memory = memory
			}
		})
		submissions = append(submissions, submissionRow)
	})

	return submissions, isDone
}
