package usaco

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-rod/rod"
)

type (
	Args struct {
		Cpid string
	}
)

// Set errors returned by library.
var (
	ErrInvalidSpecifier   = fmt.Errorf("invalid specifier data")
	errInvalidCredentials = fmt.Errorf("invalid login credentials")
)

var (
	hostURL = "http://usaco.org"

	// Browser is the headless browser to use.
	Browser *rod.Browser
)

// Parse passed in specifier string to new Args struct.
// Validates parsed args and returns error if any.
//
// List of valid specifiers can be viewed at
// TO ADD
func Parse(str string) (Args, error) {
	var (
		rxCpid = `(?P<cpid>\d+)`

		valRx = []string{
			`usaco.org\/index.php?page=viewproblem2&cpid=` + rxCpid + `$`,

			`^\s*` + rxCpid + `$`,
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
				Cpid: result["cpid"],
			}
			return arg, nil
		}
	}
	return Args{}, ErrInvalidSpecifier
}

func login(usr, passwd string) (string, error) {
	link := hostURL
	page, err := loadPage(link)
	if err != nil {
		return "", err
	}
	defer page.Close()

	// check if current user session is logged in
	if elm := page.Elements(selCSSHandle).First(); elm != nil {
		return clean(elm.Text()), nil
	}

	// otherwise, login
	page.Element(`input[name="uname"]`).Input(usr)
	page.Element(`input[name="password"]`).Input(passwd)
	page.Element(`input[value="Login"]`).Click()

	elm := page.ElementMatches(selCSSHandle, `.*`,
		selCSSFormErr, `Incorrect password`)
	if elm.Matches(selCSSFormErr) {
		return "", errInvalidCredentials
	}

	return clean(elm.Text()), nil
}

func logout() error {
	page, err := loadPage(hostURL)
	if err != nil {
		return err
	}
	defer page.Close()

	if page.HasMatches(`button`, `Logout`) {
		page.ElementMatches(`button`, `Logout`).Click()
		// page gives a notification on logout
		page.Element(`input[value="Login"]`)
	}
	return nil
}
