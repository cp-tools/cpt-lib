package atcoder

import (
	"regexp"
	"time"
)

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

	p, err := loadPage(link)
	if err != nil {
		return 0, err
	}
	defer p.Close()

	_, err = p.Race().Element(`.alert`).Handle(handleErrMsg).
		Element(`footer.footer`).Do()

	if err != nil {
		return 0, err
	}

	return p.getCountdown()
}
