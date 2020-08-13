package codeforces

import (
	"fmt"
	"regexp"
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

	// Dashboard holds details from contest dashboard.
	Dashboard struct {
		Name      string
		Problem   []Problem
		Countdown time.Duration
		// href link => description
		Material map[string]string
	}
)

// Contest registration status of current session.
const (
	RegistrationClosed    = 0
	RegistrationOpen      = 1
	RegistrationDone      = 2
	RegistrationNotExists = -1
)

// CountdownPage returns link to countdown in contest
func (arg Args) CountdownPage() (link string) {
	if arg.Class == ClassGroup {
		link = fmt.Sprintf("%v/group/%v/contest/%v/countdown",
			hostURL, arg.Group, arg.Contest)
	} else {
		link = fmt.Sprintf("%v/%v/%v/countdown",
			hostURL, arg.Class, arg.Contest)
	}
	return
}

// ContestsPage returns link to all contests page (group/gym/contest)
func (arg Args) ContestsPage() (link string) {
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

// DashboardPage returns link to dashboard of contest
func (arg Args) DashboardPage() (link string) {
	if arg.Class == ClassGroup {
		link = fmt.Sprintf("%v/group/%v/contest/%v",
			hostURL, arg.Group, arg.Contest)
	} else {
		link = fmt.Sprintf("%v/%v/%v",
			hostURL, arg.Class, arg.Contest)
	}
	return
}

// RegisterPage returns link to registration (nor virtual reg) in contest
func (arg Args) RegisterPage() (link string) {
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

	link := arg.CountdownPage()
	page, msg, err := loadPage(link)
	if err != nil {
		return 0, err
	}
	defer page.Close()

	if msg != "" {
		// there should be no notification
		return 0, fmt.Errorf(msg)
	}

	doc, _ := goquery.NewDocumentFromReader(
		strings.NewReader(page.Element("html").HTML()))

	val := doc.Find("span.countdown>span").AttrOr("title", "")
	if len(val) == 0 {
		val = doc.Find("span.countdown").Text()
	}

	var h, m, s int64
	fmt.Sscanf(val, "%d:%d:%d", &h, &m, &s)
	dur := time.Duration(h*3600+m*60+s) * time.Second
	return dur, nil
}

// GetContests extracts contest/gym/group contests data based
// on specified data in Args. Expects arg.Class to be configured
// to fetch respective contest details.
//
// Set 'count' to the maximum number of rows you want to be returned.
// Set to -1 if you want to fetch all rows of data.
func (arg Args) GetContests(count int) ([]Contest, error) {
	// MUST define Class type.
	if arg.Class != ClassGym && arg.Class != ClassGroup && arg.Class != ClassContest {
		return nil, ErrInvalidSpecifier
	}

	if count < 0 {
		// is this large enough?
		count = 1e9
	}

	link := arg.ContestsPage()
	page, msg, err := loadPage(link)
	if err != nil {
		return nil, err
	}
	defer page.Close()

	if msg != "" {
		return nil, fmt.Errorf(msg)
	}

	contests := make([]Contest, 0)
	// run till 'count' rows are parsed
	for true {
		doc, _ := goquery.NewDocumentFromReader(
			strings.NewReader(page.Element("html").HTML()))

		table := new(goquery.Selection)
		if doc.Find(".pagination span.active").AttrOr("pageindex", "1") == "1" {
			// include the upcoming contests table too
			table = doc.Find("tr[data-contestid]")
		} else {
			// exclude upcoming contests table
			table = doc.Find(".datatable").
				Eq(1).Find("tr[data-contestid]")
		}

		table.EachWithBreak(func(_ int, row *goquery.Selection) bool {
			if count == 0 {
				// Got required amount of rows. Break
				return false
			}

			// parse duration string (using ugly regex)
			parseDuration := func(str string) time.Duration {
				re := regexp.MustCompile(`(?:(\d+):)?(\d+):(\d+)`)
				val := re.FindStringSubmatch(str)
				d, _ := strconv.Atoi(val[1])
				h, _ := strconv.Atoi(val[2])
				m, _ := strconv.Atoi(val[3])
				return time.Duration(d*1440+h*60+m) * time.Minute
			}

			var contestRow Contest
			// extract contest args from html attr label
			contArg, _ := Parse(arg.Group + row.AttrOr("data-contestid", ""))
			if arg.Contest != "" && arg.Contest != contArg.Contest {
				// contest id is specified to fetch. This contest doesn't match it.
				return true
			}
			contestRow.Arg = contArg

			// the table format for contests is different from groups and gyms/contests.
			if contArg.Class == ClassContest || (contArg.Class == ClassGym && contArg.Contest != "") {
				row.Find("td").Each(func(cellIdx int, cell *goquery.Selection) {
					switch cellIdx {
					case 0:
						// remove all links from text
						cell.Find("a").Remove()
						contestRow.Name = clean(cell.Text())

					case 1:
						writers := strings.Split(clean(cell.Text()), "\n")
						if writers[0] == "" {
							// no writers are specified. Set slice to nil
							writers = nil
						}

						contestRow.Writers = writers

					case 2:
						startTime := parseTime(cell.Text())
						contestRow.StartTime = startTime

					case 3:
						duration := parseDuration(cell.Text())
						contestRow.Duration = duration

					case 5:
						cell.Find(".countdown").Remove()
						if contArg.Class == ClassGym {
							contestRow.RegStatus = RegistrationNotExists
							contestRow.RegCount = RegistrationNotExists
							description := strings.Split(clean(cell.Text()), "\n")
							contestRow.Description = description
						} else {
							// extract registration count
							cntStr := getText(cell, ".contestParticipantCountLinkMargin")
							if len(cntStr) > 1 {
								regCount, _ := strconv.Atoi(cntStr[1:])
								contestRow.RegCount = regCount
							}
							// extract registration status
							if cell.Find(".welldone").Length() != 0 {
								contestRow.RegStatus = RegistrationDone
							} else if cell.Find("a").Not("a[title]").Length() > 0 {
								contestRow.RegStatus = RegistrationOpen
							} else {
								contestRow.RegStatus = RegistrationClosed
							}
						}
					}
				})
			} else {
				row.Find("td").Each(func(cellIdx int, cell *goquery.Selection) {
					switch cellIdx {
					case 0:
						// remove all links from text
						cell.Find("a").Remove()
						contestRow.Name = clean(cell.Text())

					case 1:
						startTime := parseTime(cell.Text())
						contestRow.StartTime = startTime

					case 2:
						duration := parseDuration(cell.Text())
						contestRow.Duration = duration

					case 4:
						var description []string
						cell.Find(".small").Each(func(_ int, val *goquery.Selection) {
							description = append(description, clean(val.Text()))
						})
						contestRow.Description = description
					}
				})

				contestRow.Writers = nil
				contestRow.RegCount = RegistrationNotExists
				contestRow.RegStatus = RegistrationNotExists
			}

			contests = append(contests, contestRow)
			count--
			return true
		})

		if count == 0 {
			// got required amount of rows. Break
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

	return contests, nil
}

// GetDashboard parses and returns useful info from
// contest dashboard page.
func (arg Args) GetDashboard() (Dashboard, error) {
	if len(arg.Contest) == 0 {
		return Dashboard{}, ErrInvalidSpecifier
	}

	link := arg.DashboardPage()
	page, msg, err := loadPage(link)
	if err != nil {
		return Dashboard{}, err
	}
	defer page.Close()

	if msg != "" {
		return Dashboard{}, fmt.Errorf(msg)
	}

	doc, _ := goquery.NewDocumentFromReader(
		strings.NewReader(page.Element("html").HTML()))

	var dashboard Dashboard
	dashboard.Material = make(map[string]string)
	// extract contest name
	dashboard.Name = clean(doc.Find(".rtable th").Text())
	// extract countdown to contest end
	if true {
		str := clean(doc.Find(".countdown").Text())
		var h, m, s int
		fmt.Sscanf(str, "%d:%d:%d", &h, &m, &s)
		countdown := time.Duration(h*3600+m*60+s) * time.Second
		dashboard.Countdown = countdown
	}

	// extract problems data
	table := doc.Find(".problems tr").Has("td")
	table.Each(func(_ int, row *goquery.Selection) {
		var problemRow Problem

		// what do I do if there is an error?
		probArg, _ := Parse(hostURL + row.Find("td a").AttrOr("href", ""))
		if arg.Problem != "" && arg.Problem != probArg.Problem {
			return
		}
		problemRow.Arg = probArg

		// extract solve status
		switch row.AttrOr("class", "") {
		case "accepted-problem":
			problemRow.SolveStatus = SolveAccepted
		case "rejected-problem":
			problemRow.SolveStatus = SolveRejected
		default:
			problemRow.SolveStatus = SolveNotAttempted
		}

		row.Find("td").Each(func(cellIdx int, cell *goquery.Selection) {
			switch cellIdx {
			case 1:
				conSel := cell.Find(".notice")
				// extract time/memory limit from problem
				constraints := clean(conSel.Contents().Last().Text())
				problemRow.TimeLimit = strings.Split(constraints, ", ")[0]
				problemRow.MemoryLimit = strings.Split(constraints, ", ")[1]

				// extract input/output stream.
				if sval := getText(conSel, "div"); sval == "standard input/output" {
					problemRow.InpStream = "standard input"
					problemRow.OutStream = "standard output"
				} else {
					problemRow.InpStream = strings.Split(sval, "/")[0]
					problemRow.OutStream = strings.Split(sval, "/")[1]
				}

				name := cell.Find("a").Text()
				problemRow.Name = name

			case 3:
				solveCount := 0
				if sval := clean(cell.Text()); len(sval) > 1 {
					// remove the 'x' prefix from x123 count
					solveCount, _ = strconv.Atoi(sval[1:])
				}
				problemRow.SolveCount = solveCount
			}
		})
		dashboard.Problem = append(dashboard.Problem, problemRow)
	})

	// extract contest material
	doc.Find("#sidebar li a").Each(func(_ int, data *goquery.Selection) {
		href := data.AttrOr("href", "")
		dashboard.Material[hostURL+href] = clean(data.Text())
	})

	return dashboard, nil
}

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

	link := arg.RegisterPage()
	page, msg, err := loadPage(link)
	if err != nil {
		return nil, err
	}

	if msg != "" {
		return nil, fmt.Errorf(msg)
	}

	doc, _ := goquery.NewDocumentFromReader(
		strings.NewReader(page.Element("html").HTML()))

	registerInfo := &RegisterInfo{
		Name:  getText(doc.Selection, "h2"),
		Terms: getText(doc.Selection, ".terms"),
		Register: func() error {
			page.Element(".submit").Click()
			page.Element(`.contestList`)
			defer page.Close()
			return nil
		},
	}
	return registerInfo, nil
}
