package atcoder

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cp-tools/cpt-lib/util"

	"github.com/go-rod/rod"
)

type (
	// Args holds specifier details parsed by
	// Parse() function. All methods use this
	// at the core.
	Args struct {
		Contest string
		Problem string
	}
)

// Errors returned by library.
var (
	ErrInvalidSpecifier   = fmt.Errorf("invalid specifier data")
	errInvalidCredentials = fmt.Errorf("invalid login credentials")
)

var (
	hostURL = "https://atcoder.jp"

	// Browser is the headless browser to use.
	Browser *rod.Browser
)

// Start initiates the automated browser to use.
func Start(headless bool, userDataDir, bin string) error {
	bs, err := util.NewBrowser(headless, userDataDir, bin)
	Browser = bs

	return err
}

// loginPage returns link to login page
func loginPage() string {
	return fmt.Sprintf("%v/login", hostURL)
}

// Parse passed in specifier string to new Args struct.
// Validates parsed args and returns error if any.
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

	str = strings.TrimSpace(util.StrClean(str))
	if str == "" {
		return Args{}, nil
	}

	for _, rgx := range valRx {
		re := regexp.MustCompile(rgx)
		if re.MatchString(str) {
			// https://stackoverflow.com/a/46202939/9606036
			match := re.FindStringSubmatch(str)
			result := map[string]string{}
			for i, name := range re.SubexpNames() {
				if i != 0 && name != "" {
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
	page, msg, err := loadPage(link, selCSSFooter)
	if err != nil {
		return "", err
	}
	defer page.Close()

	if msg != "" {
		// There shouldn't be any error.
		return "", fmt.Errorf(msg)
	}

	// Check if current user is logged in.
	if !page.MustHasR(selCSSHandle, `Sign In`) {
		handle := page.MustElement(selCSSHandle).MustText()
		return util.StrClean(handle), nil
	}

	// check if username/password are valid
	if usr == "" || passwd == "" {
		return "", errInvalidCredentials
	}

	// Otherwise, login.
	page.MustElement("#username").Input(usr)
	page.MustElement("#password").Input(passwd)
	page.MustElement("#submit").MustClick().WaitInvisible()

	elm := page.Race().
		Element(selCSSHandle).
		Element(selCSSNotif + `.alert-danger`).
		MustDo()

	if elm.MustMatches(selCSSNotif) {
		return "", errInvalidCredentials
	}
	return util.StrClean(elm.MustText()), nil
}

func logout() error {
	page, msg, err := loadPage(hostURL, selCSSFooter)
	if err != nil {
		return err
	}
	defer page.Close()

	if msg != "" {
		return fmt.Errorf(msg)
	}

	if !page.MustHasR(selCSSHandle, `Sign In`) {
		// Run the logout javascript function.
		page.MustEval("form_logout.submit()")
		// Wait till logout is completed.
		page.ElementR(selCSSHandle, `Sign In`)
	}

	return nil
}
