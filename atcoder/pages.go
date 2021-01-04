package atcoder

import "fmt"

// loginPage returns link to login page
func loginPage() string {
	return fmt.Sprintf("%v/login", hostURL)
}

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
