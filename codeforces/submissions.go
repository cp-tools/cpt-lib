package codeforces

import (
	"fmt"
	"strings"
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
// Set 'count' to the maximum number of rows you want to be returned.
// Set to -1 if you want to fetch all rows of data.
func (arg Args) GetSubmissions(handle string, count int) ([]Submission, error) {
	if count < 0 {
		count = 1e9
	}

	link := arg.SubmissionsPage(handle)
	page, msg, err := loadPage(link)
	if err != nil {
		return nil, err
	}
	defer page.Close()

	if msg != "" {
		return nil, fmt.Errorf(msg)
	}

	// @todo Add support for excluding unofficial submissions

	submissions := make([]Submission, 0)
	// run till 'count' rows are parsed
	for true {
		doc, _ := goquery.NewDocumentFromReader(
			strings.NewReader(page.Element("html").HTML()))

		table := doc.Find("tr[data-submission-id]")
		table.EachWithBreak(func(_ int, row *goquery.Selection) bool {
			if count == 0 {
				// got required amount of rows. Break
				return false
			}

			var submissionRow Submission
			// extract contest args from html attr label
			subArg, _ := Parse(hostURL + getAttr(row, "td:nth-of-type(4) a", "href"))
			if arg.Problem != "" && arg.Problem != subArg.Problem {
				return true
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
			count--
			return true
		})

		if count == 0 {
			break
		}

		// navigate to next page
		if !page.HasMatches(".pagination li", "→") {
			// no more pages more left. Break
			break
		}

		// click navigation button and wait till loads
		page.ElementMatches(".pagination li", "→").Click()
		page.Element(selCSSFooter)
	}

	return submissions, nil
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
	page, msg, err := loadPage(link)
	if err != nil {
		return "", err
	}
	defer page.Close()

	if msg != "" {
		return "", fmt.Errorf(msg)
	}

	// extract source code from html body
	doc, _ := goquery.NewDocumentFromReader(
		strings.NewReader(page.Element("html").HTML()))

	source := ""
	codeBlock := doc.Find("pre#program-source-text li")
	codeBlock.Each(func(_ int, ln *goquery.Selection) {
		source += ln.Text() + "\n"
	})
	return source, nil
}
