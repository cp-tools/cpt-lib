package atcoder

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-rod/rod"
)

type (
	Args struct {
		Contest string
		Problem string
	}
)

var (
	ErrInvalidSpecifier   = fmt.Errorf("invalid specifier data")
	ErrInvalidCredentials = fmt.Errorf("invalid login credentials")
)

var (
	hostURL = "https://atcoder.jp"

	Browser *rod.Browser
)

// loginPage returns link to login page
func loginPage() string {
	return fmt.Sprintf("%v/login", hostURL)
}

func Parse(str string) (Args, error) {
	var (
		rxCont = `(?P<cont>[A-Za-z0-9-]+)`
		rxProb = `(?P<prob>[A-Za-z0-9_]+)`

		valRx = []string{
			`atcoder.jp\/contests\/` + rxCont + `$`,
			`atcoder.jp\/contests\/` + rxCont + `\/tasks\/` + rxProb + `$`,
		}
	)

	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return Args{}, nil
	}

	for _, rgx := range valRx {
		re := regexp.MustCompile(rgx)
		if re.MatchString(str) {
			// attrib : stackoverflow.com/a/9606036
			match := re.FindStringSubmatch(str)
			result := map[string]string{}
			for i, name := range re.SubexpNames() {
				if i != 0 && len(name) > 0 {
					result[name] = match[i]
				}
			}
			arg := Args{
				Contest: result["cont"],
				Problem: result["prob"],
			}
			return arg, nil
		}
	}
	return Args{}, ErrInvalidSpecifier
}

func login(usr, passwd string) (string, error) {
	link := loginPage()
	page, err := Browser.PageE(link)
	if err != nil {
		return "", err
	}

	page.WaitLoad()
	if msg := cE(page); msg != "" {
		return "", fmt.Errorf(msg)
	}

	page.Element(".navbar-right")
	// check if current user sesion is logged in
	fmt.Println(findHandle(page))
	if handle := findHandle(page); handle != "" {
		return handle, nil
	}

	// otherwise, login
	page.Element("#username").Input(usr)
	page.Element("#password").Input(passwd)
	page.Element("#submit").Click()

	// race 2 selectors, wait till one resolves
	errSelector := ".alert-danger"
	el := page.Element(errSelector, ".alert-success")
	if el.Matches(errSelector) {
		return "", ErrInvalidCredentials
	}

	handle := findHandle(page)
	return handle, nil
}
