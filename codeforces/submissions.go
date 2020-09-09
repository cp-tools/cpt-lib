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
func (arg Args) SubmissionsPage(handle string) (link string) {
	// contest specified
	if len(arg.Contest) != 0 {
		if arg.Class == ClassGroup {
			// does this even work?!
			if len(handle) == 0 {
				link = fmt.Sprintf("%v/group/%v/contest/%v/my",
					hostURL, arg.Group, arg.Contest)
			} else {
				link = fmt.Sprintf("%v/submissions/%v/group/%v/contest/%v",
					hostURL, handle, arg.Group, arg.Contest)
			}
		} else {
			if len(handle) == 0 {
				link = fmt.Sprintf("%v/%v/%v/my",
					hostURL, arg.Class, arg.Contest)
			} else {
				link = fmt.Sprintf("%v/submissions/%v/%v/%v",
					hostURL, handle, arg.Class, arg.Contest)
			}
		}
	} else {
		if len(handle) == 0 {
			// I think this is a bad idea....
			handle, _ = login("", "")
		}
		link = fmt.Sprintf("%v/submissions/%v",
			hostURL, handle)
	}
	return
}

// SourceCodePage returns link to solution submission
func (sub Submission) SourceCodePage() (link string) {
	if sub.Arg.Class == ClassGroup {
		link = fmt.Sprintf("%v/group/%v/contest/%v/submission/%v",
			hostURL, sub.Arg.Group, sub.Arg.Contest, sub.ID)
	} else {
		link = fmt.Sprintf("%v/%v/%v/submission/%v",
			hostURL, sub.Arg.Class, sub.Arg.Contest, sub.ID)
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
	if pageCount < 0 {
		pageCount = 1e9
	}

	link := arg.SubmissionsPage(handle)
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
				page.MustElementMatches(".pagination li", "→").MustClick()
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
	if len(sub.Arg.Contest) == 0 || len(sub.ID) == 0 {
		return "", ErrInvalidSpecifier
	}

	link := sub.SourceCodePage()
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
