package codeforces

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type (
	// Contest holds details from contest row
	// from contests table.
	Contest struct {
		Name        string
		Writers     []string
		StartTime   time.Time
		Duration    time.Duration
		RegCount    int
		RegStatus   int
		Description []string
		Arg         Args
	}

	// RegisterInfo holds data pertaining to contest
	// registration along with a callback function to
	// register in the said contest.
	RegisterInfo struct {
		Name     string
		Terms    string
		Register func() error
	}

	/*
		// Dashboard holds details from contest dashboard
		Dashboard struct {}
	*/
)

// Contest registration status of current session.
const (
	RegistrationClosed    = 0
	RegistrationOpen      = 1
	RegistrationDone      = 2
	RegistrationNotExists = -1
)

func (arg Args) countdownPage() (link string) {
	if arg.Class == ClassGroup {
		link = fmt.Sprintf("%v/group/%v/contest/%v/countdown",
			hostURL, arg.Group, arg.Contest)
	} else {
		link = fmt.Sprintf("%v/%v/%v/countdown",
			hostURL, arg.Class, arg.Contest)
	}
	return
}

func (arg Args) contestsPage() (link string) {
	if arg.Class == ClassGroup {
		// details of individual contest can't be parsed.
		// fallback to parsing all contests in group.
		link = fmt.Sprintf("%v/group/%v/contests?complete=true",
			hostURL, arg.Group)
	} else if len(arg.Contest) != 0 {
		link = fmt.Sprintf("%v/contests/%v",
			hostURL, arg.Contest)
	} else {
		link = fmt.Sprintf("%v/%vs?complete=true",
			hostURL, arg.Class)
	}
	return
}

func (arg Args) registerPage() (link string) {
	// gyms/groups don't support registration, do they!?
	link = fmt.Sprintf("%v/contestRegistration/%v",
		hostURL, arg.Contest)
	return
}

// GetCountdown parses and returns duration type for countdown
// in specified contest to end. If countdown has already ended,
// returns 0. Extracts data from .../contest/<contest>/countdown.
func (arg Args) GetCountdown() (time.Duration, error) {
	if len(arg.Contest) == 0 {
		return 0, ErrInvalidSpecifier
	}

	link := arg.countdownPage()
	resp, err := SessCln.Get(link)
	if err != nil {
		return 0, err
	}
	body, msg := parseResp(resp)
	if len(msg) != 0 {
		return 0, fmt.Errorf(msg)
	}

	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	val := doc.Find("span.countdown").Text()

	var h, m, s int64
	fmt.Sscanf(val, "%d:%d:%d", &h, &m, &s)
	dur := time.Duration(h*3600+m*60+s) * time.Second
	return dur, nil
}

// GetContests extracts contest/gym/group contests data based
// on specified data in Args. Expects arg.Class to be configured
// to fetch respective contest details.
//
// Set 'omitFinishedContests' to true to exclude finished contests.
func (arg Args) GetContests(omitFinishedContests bool) ([]Contest, error) {
	// MUST define Class type.
	if arg.Class != ClassGym && arg.Class != ClassGroup && arg.Class != ClassContest {
		return nil, ErrInvalidSpecifier
	}

	link := arg.contestsPage()
	resp, err := SessCln.Get(link)
	if err != nil {
		return nil, err
	}
	body, msg := parseResp(resp)
	if len(msg) != 0 {
		return nil, fmt.Errorf(msg)
	}

	var contests []Contest
	pages := findPagination(body)
	for c := 1; c <= pages; c++ {
		cLink := fmt.Sprintf("%v/page/%d", link, c)
		resp, err = SessCln.Get(cLink)
		if err != nil {
			return nil, err
		}
		body, msg = parseResp(resp)
		if len(msg) != 0 {
			return nil, fmt.Errorf(msg)
		}

		isOver := false
		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
		table := doc.Find("tr[data-contestid]")
		table.EachWithBreak(func(_ int, cont *goquery.Selection) bool {
			// extract contest args from html attr label
			contArg, _ := Parse(arg.Group + " " + cont.AttrOr("tr[data-contestid]", ""))
			contArg.setContestClass()

			// remove links from contest name
			name := cont.Find("td:nth-of-type(1)")
			name.Find("a").Remove()

			if len(arg.Contest) != 0 && strings.EqualFold(arg.Contest, contArg.Contest) {
				// skip current contest data
				return true
			}

			// extract duration from contest length
			parseDur := func(str string) time.Duration {
				d, h, m := 0, 0, 0
				// format - days:hours:minutes
				_, err := fmt.Sscanf(str, "%d:%d:%d", &d, &h, &m)
				if err != nil {
					d, h, m = 0, 0, 0
					// format - hours:minutes
					fmt.Sscanf(str, "%d:%d", &h, &m)
				}
				dur := time.Duration(d*1440+h*60+m) * time.Minute
				return dur
			}

			// handle different table formats
			if arg.Class == ClassGroup || (arg.Class == ClassGym && len(arg.Contest) == 0) {
				startTime := parseTime(getText(cont, "td:nth-of-type(2)"))
				dur := parseDur(getText(cont, "td:nth-of-type(3)"))

				if omitFinishedContests == true && time.Now().After(startTime.Add(dur)) {
					// break out of loop
					isOver = true
					return false
				}

				var description []string
				cont.Find("td:nth-of-type(5) .small").Each(func(_ int, desc *goquery.Selection) {
					description = append(description, clean(desc.Text()))
				})

				contests = append(contests, Contest{
					Name:        clean(name.Text()),
					Writers:     nil,
					StartTime:   startTime,
					Duration:    dur,
					RegCount:    RegistrationNotExists,
					RegStatus:   RegistrationNotExists,
					Description: description,
					Arg:         contArg,
				})
			} else {
				startTime := parseTime(getText(cont, "td:nth-of-type(3)"))
				dur := parseDur(getText(cont, "td:nth-of-type(4)"))

				if omitFinishedContests == true && time.Now().After(startTime.Add(dur)) {
					// break out of loop
					isOver = true
					return false
				}

				writers := strings.Split(getText(cont, "td:nth-of-type(2)"), "\n")

				// find registration state in contest
				status := cont.Find("td:nth-of-type(6)")
				status.Find(".countdown").Remove()
				var regStatus, regCount int
				if arg.Class == ClassGym {
					regStatus = RegistrationNotExists
					regCount = RegistrationNotExists
				} else {
					// extract registration count
					cntStr := getText(cont, ".contestParticipantCountLinkMargin")
					regCount, _ = strconv.Atoi(cntStr)
					// extract registration status
					if status.Find(".welldone").Length() != 0 {
						regStatus = RegistrationDone
					} else if status.Find("a").Not("a[title]").Length() > 0 {
						regStatus = RegistrationOpen
					} else {
						regStatus = RegistrationClosed
					}
				}

				contests = append(contests, Contest{
					Name:        clean(name.Text()),
					Writers:     writers,
					StartTime:   startTime,
					Duration:    dur,
					RegCount:    regCount,
					RegStatus:   regStatus,
					Description: nil,
					Arg:         contArg,
				})
			}
			return true
		})
		if isOver == true {
			break
		}
	}
	return contests, nil
}

// @todo Add GetDashboard functionality
// @body Extract and return problems in contest, time to
// @body contest end (if any), and announcements.

/*func (arg Args) GetDashboard() (Dashboard, error) {}*/

// RegisterForContest parses and returns registration terms
// of contest specified in args.
//
// Provides callback method to register current user session
// in contest. If registration was successful, returns nil error.
func (arg Args) RegisterForContest() (*RegisterInfo, error) {
	// ONLY contests support registration
	if arg.Class != ClassContest || len(arg.Contest) == 0 {
		return nil, ErrInvalidSpecifier
	}

	link := arg.registerPage()
	resp, err := SessCln.Get(link)
	if err != nil {
		return nil, err
	}
	body, msg := parseResp(resp)
	if len(msg) != 0 {
		return nil, fmt.Errorf(msg)
	}

	// hidden form data
	csrf := findCsrf(body)
	ftaa := genRandomString(18)
	bfaa := genRandomString(32)

	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	registerInfo := &RegisterInfo{
		Name:  getText(doc.Selection, "h2"),
		Terms: getText(doc.Selection, ".terms"),
		Register: func() error {
			_, err := SessCln.PostForm(link, url.Values{
				"csrf_token": {csrf},
				"ftaa":       {ftaa},
				"bfaa":       {bfaa},
				"action":     {"formSubmitted"},
				"backUrl":    {""},
				"takePartAs": {"personal"},
				"_tta":       {"176"},
			})
			return err
		},
	}
	return registerInfo, nil
}
