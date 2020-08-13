package atcoder

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-rod/rod"
)

type (
	// Args holds specifier details parsed by
	// Parse() function.
	Args struct {
		Contest string
		Problem string
	}
)

// Set errors returned by library.
var (
	ErrInvalidSpecifier   = fmt.Errorf("invalid specifier data")
	ErrInvalidCredentials = fmt.Errorf("invalid login credentials")
)

var (
	hostURL = "https://atcoder.jp"

	// Browser is the headless browser to use.
	Browser *rod.Browser
)

// loginPage returns link to login page
func loginPage() string {
	return fmt.Sprintf("%v/login", hostURL)
}

// Parse passed in specifier string to new Args struct.
// Validates parsed args and returns error if any.
//
// List of valid specifiers can be viewed at:
// TO ADD **
func Parse(str string) (Args, error) {
	var (
		rxCont = `(?P<cont>[A-Za-z0-9-]+)`
		rxProb = `(?P<prob>[A-Za-z0-9_]+)`

		valRx = []string{
			`atcoder.jp\/contests\/` + rxCont + `$`,
			`atcoder.jp\/contests\/` + rxCont + `\/tasks\/` + rxProb + `$`,

			`^` + rxCont + `$`,
			`^` + rxCont + `\s+` + rxProb + `$`,
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
	page, msg, err := loadPage(link)
	if err != nil {
		return "", err
	}
	defer page.Close()

	if msg != "" {
		// there shouldn't be any notification
		return "", fmt.Errorf(msg)
	}

	// check if current user is logged in
	if !page.HasMatches(`.nav>li>a`, `Sign In`) {
		handle := page.Element(selCSSHandle).Text()
		return clean(handle), nil
	}

	// otherwise login
	page.Element("#username").Input(usr)
	page.Element("#password").Input(passwd)
	page.Element("#submit").Click()

	elm := page.Element(selCSSHandle, selCSSNotif)
	if elm.Matches(selCSSNotif) {
		return "", ErrInvalidCredentials
	}
	return clean(elm.Text()), nil
}

func logout() error {
	page, msg, err := loadPage(hostURL)
	if err != nil {
		return err
	}
	defer page.Close()

	if msg != "" {
		return fmt.Errorf(msg)
	}

	if page.HasMatches(`a`, `Sign Out`) {
		page.ElementMatches(`a`, `Sign Out`).Click()
		// page gives notification on logout
		page.Element(selCSSNotif)
	}
	return nil
}

// trigger build
