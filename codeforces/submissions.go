package codeforces

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type (
	// Submission holds submission data
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

// SubmissionsPage returns link to user submissions
func (arg Args) SubmissionsPage(handle string) (link string, err error) {
	// contest specified
	if arg.Contest != "" {
		if arg.Class == ClassGroup {
			// groups are not supported.
			return "", ErrInvalidSpecifier
		}

		if arg.Class == "" || arg.Contest == "" {
			return "", ErrInvalidSpecifier
		}

		if handle == "" {
			link = fmt.Sprintf("%v/%v/%v/my",
				hostURL, arg.Class, arg.Contest)
		} else {
			link = fmt.Sprintf("%v/submissions/%v/%v/%v",
				hostURL, handle, arg.Class, arg.Contest)
		}

	} else {
		if handle == "" {
			tmpHandle, err := login("", "")
			if err != nil {
				return "", err
			}
			handle = tmpHandle
		}

		link = fmt.Sprintf("%v/submissions/%v",
			hostURL, handle)
	}
	return
}

// SourceCodePage returns link to solution submission
func (sub Submission) SourceCodePage() (link string, err error) {
	arg := sub.Arg // becomes too long to type otherwise...

	if (arg.Class != ClassContest && arg.Class != ClassGym && arg.Class != ClassGroup) ||
		arg.Contest == "" || sub.ID == "" {
		return "", ErrInvalidSpecifier
	}

	if arg.Class == ClassGroup {
		link = fmt.Sprintf("%v/group/%v/contest/%v/submission/%v",
			hostURL, arg.Group, arg.Contest, sub.ID)
	} else {
		link = fmt.Sprintf("%v/%v/%v/submission/%v",
			hostURL, arg.Class, arg.Contest, sub.ID)
	}
	return
}

// GetSubmissions parses and returns all submissions data in specified args
// of given user. Fetches details of all submissions of handle if args is nil.
//
// If handle is not set, fetches submissions of currently active user session.
// Due to a bug on codeforces, submissions in groups are not supported.
//
// Set 'pageCount' to the maximum number of pages of rows you want to be scraped.
// If 'pageCount' > 1, channel will not wait until all verdicts are declared, and will
// close once verdicts from all specified pages are extracted.
func (arg Args) GetSubmissions(handle string, pageCount int) (<-chan []Submission, error) {
	if pageCount <= 0 {
		return nil, nil
	}

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
	chanSubmissions := make(chan []Submission, 500)
	go func() {
		defer page.Close()
		defer close(chanSubmissions)

		// parse submissions from current page
		parseFunc := func() ([]Submission, bool) {
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

		if pageCount == 1 {
			// loop till 'isDone' is true
			for true {
				submissions, isDone := parseFunc()
				chanSubmissions <- submissions
				if isDone == true {
					break
				}
				time.Sleep(time.Millisecond * 500)
			}
		} else {
			// iterate till no more valid required pages left
			for ; pageCount > 0; pageCount-- {
				submissions, _ := parseFunc()
				chanSubmissions <- submissions

				if !page.MustHasMatches(".pagination li", "→") || pageCount == 0 {
					// no more pages to parse
					break
				}
				// click navigation button and wait till loads
				page.MustElementR(".pagination li", "→").MustClick()
				page.Element(`tr[data-submission-id]`)
			}
		}
	}()
	return chanSubmissions, nil
}

// GetSourceCode parses and returns source code of submission
// as specified in the method. Has an auto sleep cycle of 4 seconds
// to handle http error "Too Many Requests".
//
// Due to a bug on codeforces, groups are not supported.
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
	source := page.MustEval(`document.querySelector(
		"#program-source-text").innerText`).String()
	return source, nil
}
