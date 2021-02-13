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

	// Dashboard holds details from contest dashboard.
	Dashboard struct {
		Name      string
		Problem   []Problem
		Countdown time.Duration
		// href link => description
		Material map[string]string
	}
)

// Contest registration status.
const (
	RegistrationClosed    = 0
	RegistrationOpen      = 1
	RegistrationDone      = 2
	RegistrationNotExists = -1
)

func (p *page) getCountdown() (time.Duration, error) {
	pd := p.parse()

	countdownStr := pd.Find(`span.countdown>span`).AttrOr(`title`, "")
	if countdownStr == "" {
		countdownStr = pd.Find(`span.countdown`).Text()
	}

	var h, m, s int64
	fmt.Sscanf(countdownStr, "%d:%d:%d", &h, &m, &s)
	dur := time.Duration(h*3600+m*60+s) * time.Second
	return dur, nil
}

// GetCountdown returns the time before the given contest begins.
// If contest has already started, returns 0.
//
// Use this function instead of GetContests to get countdown,
// as it supports returning countdown of virtual contests too.
func (arg Args) GetCountdown() (time.Duration, error) {
	// chan has not been implemented here since,
	// countdown is updated on reload,
	// and is not websocket based.

	link, err := arg.CountdownPage()
	if err != nil {
		return 0, err
	}

	p, err := loadPage(link)
	if err != nil {
		return 0, err
	}
	defer p.Close()

	if _, err := p.Race().Element(`#jGrowl .message`).Handle(handleErrMsg).
		Element(`#footer`).Do(); err != nil {
		return 0, err
	}

	return p.getCountdown()
}

func (p *page) getContests(arg Args) ([]Contest, error) {
	pd := p.parse()

	contests := make([]Contest, 0)

	contestTableRows := pd.Find(`tr[data-contestid]`)
	contestTableRows.Each(func(_ int, row *goquery.Selection) {
		// parse duration string (using ugly regex)
		parseDuration := func(str string) time.Duration {
			re := regexp.MustCompile(`(?:(\d+):)?(\d+):(\d+)`)
			val := re.FindStringSubmatch(str)
			d, _ := strconv.Atoi(val[1])
			h, _ := strconv.Atoi(val[2])
			m, _ := strconv.Atoi(val[3])
			return time.Duration(d*1440+h*60+m) * time.Minute
		}

		// Data of current contest is stored in this.
		var contest Contest

		contest.Arg, _ = Parse(arg.Group + row.AttrOr(`data-contestid`, ``))
		if arg.Contest != "" && arg.Contest != contest.Arg.Contest {
			// contest id is specified to fetch.
			// This contest doesn't match it.
			return
		}

		// the table format for contests is different from groups and gyms/contests.
		if (contest.Arg.Class == ClassContest) ||
			(arg.Class == ClassGym && arg.Contest != "") {
			row.Find(`td`).Each(func(cellIndex int, cell *goquery.Selection) {
				switch cellIndex {
				case 0:
					cell.Find(`a`).Remove()
					contest.Name = clean(cell.Text())

				case 1:
					writers := strings.Split(clean(cell.Text()), "\n")
					if writers[0] == "" {
						// no writers are specified. Set slice to nil
						writers = nil
					}

					contest.Writers = writers

				case 2:
					contest.StartTime = parseTime(cell.Text())

				case 3:
					contest.Duration = parseDuration(clean(cell.Text()))

				case 5:
					cell.Find(`.countdown`).Remove()
					if contest.Arg.Class == ClassGym {
						contest.RegStatus = RegistrationNotExists
						contest.RegCount = RegistrationNotExists
						contest.Description = strings.Split(clean(cell.Text()), "\n")
					} else {
						// extract registration count
						registrationCountStr := cell.Find(`.contestParticipantCountLinkMargin`).Text()
						registrationCountStr = clean(registrationCountStr)

						if len(registrationCountStr) > 1 {
							regCount, _ := strconv.Atoi(registrationCountStr[1:])
							contest.RegCount = regCount
						}
						// extract registration status
						if cell.Find(`.welldone`).Length() != 0 {
							contest.RegStatus = RegistrationDone
						} else if cell.Find(`a`).Not(`a[title]`).Length() > 0 {
							contest.RegStatus = RegistrationOpen
						} else {
							contest.RegStatus = RegistrationClosed
						}
					}
				}
			})
		} else {
			row.Find("td").Each(func(cellIndex int, cell *goquery.Selection) {
				switch cellIndex {
				case 0:
					// remove all links from text
					cell.Find(`a`).Remove()
					contest.Name = clean(cell.Text())

				case 1:
					contest.StartTime = parseTime(clean(cell.Text()))

				case 2:
					contest.Duration = parseDuration(clean(cell.Text()))

				case 4:
					var description []string
					cell.Find(`.small`).Each(func(_ int, val *goquery.Selection) {
						description = append(description, clean(val.Text()))
					})
					contest.Description = description
				}
			})

			contest.Writers = nil
			contest.RegStatus = RegistrationNotExists
			contest.RegCount = RegistrationNotExists
		}

		contests = append(contests, contest)
	})
	return contests, nil
}

// GetContests returns metadata of the given contest(s).
//
// Set 'pageCount' to the maximum number of pages to parse.
// Each page consists of 100 rows of data, except the first page,
// which may contain additional upcoming contests data.
func (arg Args) GetContests(pageCount uint) (<-chan []Contest, error) {
	link, err := arg.ContestsPage()
	if err != nil {
		return nil, err
	}

	p, err := loadPage(link)
	if err != nil {
		return nil, err
	}

	if _, err := p.Race().Element(`#jGrowl .message`).Handle(handleErrMsg).
		Element(`tr[data-contestid]`).Do(); err != nil {
		p.Close()
		return nil, err
	}
	// Wait till alls rows are loaded.
	p.WaitLoad()

	chanContests := make(chan []Contest)
	go func() {
		defer p.Close()
		defer close(chanContests)

		for ; pageCount > 0; pageCount-- {
			// Ignore error, write whatever is parsed.
			contests, _ := p.getContests(arg)
			chanContests <- contests

			if !p.MustHasR(`.pagination li>a`, `→`) || pageCount == 1 {
				// All pages parsed.
				break
			}

			// Move to the next page (click the next button).
			p.MustElementR(`.pagination li>a`, `→`).MustClick().WaitInvisible()
			p.WaitLoad()

			// Remove upcoming contests table.
			if arg.Class == ClassContest {
				p.MustElements(`.contestList .datatable`).First().Remove()
			}
		}
	}()

	return chanContests, nil
}

func (p *page) getDashboard(arg Args) (Dashboard, error) {
	pd := p.parse()

	// Dashboard data is stored to this.
	var dashboard Dashboard

	dashboard.Name = pd.Find(".rtable th").Text()

	// Extract countdown to contest end.
	if countdownStr := pd.Find(`span.countdown>span`).AttrOr("title", ""); true {
		if countdownStr == "" {
			countdownStr = pd.Find(`span.countdown`).Text()
		}

		var h, m, s int
		fmt.Sscanf(countdownStr, "%d:%d:%d", &h, &m, &s)
		dashboard.Countdown = time.Duration(h*3600+m*60+s) * time.Second
	}

	problemsTable := pd.Find(`.problems tr`).Has(`td`)
	problemsTable.Each(func(_ int, row *goquery.Selection) {
		// Problem data is stored to this.
		var problem Problem

		problem.Arg, _ = Parse(hostURL + row.Find(`td a`).AttrOr(`href`, ``))
		if arg.Problem != "" && arg.Problem != problem.Arg.Problem {
			return
		}

		row.Find(`td`).Each(func(cellIndex int, cell *goquery.Selection) {
			switch cellIndex {
			case 1:
				noticeSel := cell.Find(`.notice`)
				// Extract time/memory limit data.
				constraints := clean(noticeSel.Contents().Last().Text())
				problem.TimeLimit = strings.Split(constraints, ", ")[0]
				problem.MemoryLimit = strings.Split(constraints, ", ")[1]

				// Extract input/output stream data.
				if streamStr := clean(noticeSel.Find(`div`).Text()); streamStr == "standard input/output" {
					problem.InpStream = "standard input"
					problem.OutStream = "standard output"
				} else {
					problem.InpStream = clean(strings.Split(streamStr, "/")[0])
					problem.OutStream = clean(strings.Split(streamStr, "/")[1])
				}

				problem.Name = clean(cell.Find("a").Text())

			case 3:
				if solveStr := clean(cell.Text()); len(solveStr) > 1 {
					// Remove the 'x' prefix from the string.
					problem.SolveCount, _ = strconv.Atoi(solveStr[1:])
				}
			}
		})

		// Extract solve status.
		switch row.AttrOr(`class`, ``) {
		case "accepted-problem":
			problem.SolveStatus = SolveAccepted
		case "rejected-problem":
			problem.SolveStatus = SolveRejected
		default:
			problem.SolveStatus = SolveNotAttempted
		}

		dashboard.Problem = append(dashboard.Problem, problem)
	})

	// Create map to hold material links.
	dashboard.Material = make(map[string]string)
	pd.Find(`#sidebar li a`).Each(func(_ int, sel *goquery.Selection) {
		dashboard.Material[hostURL+sel.AttrOr(`href`, ``)] = sel.Text()
	})

	return dashboard, nil
}

// GetDashboard returns in depth contest metadata from
// the contest dashboard page.
//
// Data returned by this function is user session specific,
// as user interaction in the contest is parsed and returned.
func (arg Args) GetDashboard() (Dashboard, error) {
	link, err := arg.DashboardPage()
	if err != nil {
		return Dashboard{}, err
	}

	p, err := loadPage(link)
	if err != nil {
		return Dashboard{}, err
	}
	defer p.Close()

	if _, err := p.Race().Element(`#jGrowl .message`).Handle(handleErrMsg).
		Element(`#footer`).Do(); err != nil {
		return Dashboard{}, err
	}

	return p.getDashboard(arg)
}
