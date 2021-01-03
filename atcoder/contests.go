package atcoder

import (
	"fmt"
	"regexp"
	"time"
)

// DashboardPage returns link to dashboard of contest
func (arg Args) DashboardPage() (link string, err error) {
	if arg.Contest == "" {
		return "", ErrInvalidSpecifier
	}

	link = fmt.Sprintf("%v/contests/%v", hostURL, arg.Contest)
	return
}

// VirtualPage returns link to virtual contest tab.
func (arg Args) VirtualPage() (link string, err error) {
	if arg.Contest == "" {
		return "", ErrInvalidSpecifier
	}

	dashboardLink, _ := arg.DashboardPage()
	link = fmt.Sprintf("%v/virtual", dashboardLink)
	return
}

func (p *page) getCountdown() (time.Duration, error) {
	// First check for virtual countdown.
	if match := regexp.MustCompile(`var virtualStartTime = moment\("(.*)"\)`).
		FindStringSubmatch(p.MustElement("html").MustHTML()); len(match) == 2 {
		startTime, err := time.Parse(time.RFC3339Nano, match[1])
		if err != nil {
			return 0, err
		}

		dur := time.Until(startTime).Truncate(time.Second)
		if dur < 0 { // Virtual already started; No countdown
			dur = 0
		}
		return dur, nil
	}

	// Parse actual contest start time.
	startTimeStr := p.MustEval("startTime._d.toISOString()").String()
	startTime, err := time.Parse(time.RFC3339Nano, startTimeStr)
	if err != nil {
		return 0, err
	}

	dur := time.Until(startTime).Truncate(time.Second)
	if dur < 0 { // Contest already started; No countdown
		dur = 0
	}
	return dur, nil
}

// GetCountdown ...
func (arg Args) GetCountdown() (time.Duration, error) {
	link, err := arg.VirtualPage()
	if err != nil {
		return 0, err
	}

	p, msg, err := loadPage(link, selCSSFooter)
	if err != nil {
		return 0, err
	}
	defer p.Close()

	if msg != "" {
		// there should be no notification
		return 0, fmt.Errorf(msg)
	}

	return p.getCountdown()
}
