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
			handle, _ = Login("", "")
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
//
// Due to a bug on codeforces, submissions in groups are not supported.
func (arg Args) GetSubmissions(handle string) ([]Submission, error) {
	link := arg.SubmissionsPage(handle)
	page, err := Browser.PageE(link)
	if err != nil {
		return nil, err
	}

	page.WaitLoad()
	if msg := cE(page); msg != "" {
		return nil, fmt.Errorf(msg)
	}
	body := page.Element("html").HTML()

	// @todo Add support for excluding unofficial submissions

	var submissions []Submission
	pages := findPagination(page)
	for c := 1; c <= pages; c++ {
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(body))
		doc.Find("tr[data-submission-id]").Each(func(_ int, sub *goquery.Selection) {
			var newSubmission Submission

			subArg, _ := Parse(hostURL + getAttr(sub, "td:nth-of-type(4) a", "href"))
			if arg.Problem != "" && arg.Problem != subArg.Problem {
				return
			}
			newSubmission.Arg = subArg

			sub.Find("td").Each(func(cellIdx int, cell *goquery.Selection) {
				switch cellIdx {
				case 0:
					id := clean(cell.Text())
					newSubmission.ID = id

				case 1:
					when := parseTime(clean(cell.Text()))
					newSubmission.When = when

				case 2:
					who := clean(cell.Text())
					newSubmission.Who = who

				case 3:
					problem := clean(cell.Text())
					newSubmission.Problem = problem

				case 4:
					language := clean(cell.Text())
					newSubmission.Language = language

				case 5:
					isJudging := cell.AttrOr("waiting", "") == "true"
					newSubmission.IsJudging = isJudging

					verdict := clean(cell.Text())
					newSubmission.Verdict = verdict

				case 6:
					time := clean(cell.Text())
					newSubmission.Time = time

				case 7:
					memory := clean(cell.Text())
					newSubmission.Memory = memory
				}
			})
			submissions = append(submissions, newSubmission)
		})

		if c+1 <= pages {
			cLink := fmt.Sprintf("%v/page/%d", link, c+1)
			page, err := Browser.PageE(cLink)
			if err != nil {
				return submissions, err
			}

			page.WaitLoad()
			if msg := cE(page); msg != "" {
				return nil, fmt.Errorf(msg)
			}
			body = page.Element("html").HTML()
		}
	}
	return submissions, nil
}

// GetSourceCode parses and returns source code of submission
// as specified in the method. Has an auto sleep cycle of 4 seconds
// to handle http error "Too Many Requests".
//
// Due to a bug on codeforces, groups are not supported.
func (sub Submission) GetSourceCode() (string, error) {
	var source string
	if len(sub.Arg.Contest) == 0 || len(sub.ID) == 0 {
		return "", ErrInvalidSpecifier
	}
	link := sub.SourceCodePage()
	page, err := Browser.PageE(link)
	if err != nil {
		return "", err
	}

	page.WaitLoad()
	if msg := cE(page); msg != "" {
		return "", fmt.Errorf(msg)
	}
	body := page.Element("html").HTML()

	// extract source code from html body
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(body))
	doc.Find("pre#program-source-text li").Each(func(_ int, ln *goquery.Selection) {
		source += ln.Text() + "\n"
	})
	return clean(source), nil
}
