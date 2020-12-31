package atcoder

import (
	"fmt"
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

func (arg Args) CountdownPage(isVC bool) (link string, err error) {
	if arg.Contest == "" {
		return "", ErrInvalidSpecifier
	}

	dashboardLink, _ := arg.DashboardPage()
	if isVC == true {
		link = fmt.Sprintf("%v/virtual", dashboardLink)
		return
	}

	link = dashboardLink
	return
}

func (arg Args) GetCountdown(isVC bool) (time.Duration, error) {
	if arg.Contest == "" {
		return 0, ErrInvalidSpecifier
	}

	link, err := arg.CountdownPage(isVC)
	if err != nil {
		return 0, err
	}

	page, msg, err := loadPage(link, selCSSFooter)
	if err != nil {
		return 0, err
	}
	defer page.Close()

	if msg != "" {
		// there should be no notification
		return 0, fmt.Errorf(msg)
	}

	fmt.Println(page.MustEval("startTime._d.toISOString()").String())
	return 0, nil
}
