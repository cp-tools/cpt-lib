package codeforces

import (
	"bytes"
	"fmt"
	"net/http"
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

func (arg Args) submissionsPage(handle string) (link string) {
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
		// I think this is a bad idea....
		handle, _ := Login("", "")
		link = fmt.Sprintf("%v/submissions/%v",
			hostURL, handle)
	}
	return
}

func (sub Submission) sourceCodePage() (link string) {
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
	link := arg.submissionsPage(handle)
	resp, err := SessCln.Get(link)
	if err != nil {
		return nil, err
	}
	body, msg := parseResp(resp)
	if len(msg) != 0 {
		// shouldn't return any error on success
		return nil, fmt.Errorf(msg)
	}

	// @todo Add support for excluding unofficial submissions

	var submissions []Submission
	pages := findPagination(body)
	for c := 1; c <= pages; c++ {
		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
		table := doc.Find("tr[data-submission-id]")
		table.Each(func(_ int, sub *goquery.Selection) {
			probLk := hostURL + getAttr(sub, "td:nth-of-type(4) a", "href")
			subArg, _ := Parse(probLk)

			if len(arg.Problem) == 0 || arg.Problem == subArg.Problem {
				// parse various details
				isJudging := getAttr(sub, "td:nth-of-type(6)", "waiting") == "true"
				when := parseTime(getText(sub, "td:nth-of-type(2)"))

				submissions = append(submissions, Submission{
					ID:       getText(sub, "td:nth-of-type(1)"),
					When:     when,
					Who:      getText(sub, "td:nth-of-type(3)"),
					Problem:  getText(sub, "td:nth-of-type(4)"),
					Language: getText(sub, "td:nth-of-type(5)"),
					Verdict:  getText(sub, "td:nth-of-type(6)"),
					Time:     getText(sub, "td:nth-of-type(7)"),
					Memory:   getText(sub, "td:nth-of-type(8)"),

					IsJudging: isJudging,
					Arg:       subArg,
				})
			}
		})
		if c+1 <= pages {
			cLink := fmt.Sprintf("%v/page/%d", link, c+1)
			resp, err = SessCln.Get(cLink)
			if err != nil {
				return nil, err
			}
			body, msg = parseResp(resp)
			if len(msg) != 0 {
				return nil, fmt.Errorf(msg)
			}
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

	link := sub.sourceCodePage()
FETCH:
	resp, err := SessCln.Get(link)
	if err != nil {
		return source, err
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		// sleep for 4 seconds before resending request
		time.Sleep(time.Second * 4)
		goto FETCH
	}
	body, msg := parseResp(resp)
	if len(msg) != 0 {
		// msg should be length 0 if success
		return source, fmt.Errorf(msg)
	}
	// extract source code from html body
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	source = doc.Find("pre#program-source-text").Text()
	return source, nil
}
