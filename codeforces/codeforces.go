package codeforces

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-rod/rod"
)

type (
	// Args holds specifier details parsed by
	// Parse() function. Provides methods to
	// access value of specifiers (read-only).
	Args struct {
		Contest string
		Problem string
		Class   string
		Group   string
	}
)

// Class type of contest.
const (
	ClassContest = "contest"
	ClassGroup   = "group"
	ClassGym     = "gym"
)

// Set errors returned by library.
var (
	ErrInvalidSpecifier   = fmt.Errorf("invalid specifier data")
	ErrInvalidCredentials = fmt.Errorf("invalid login credentials")
)

var (
	hostURL = "https://codeforces.com"
	/*
		// SessCln should be set to desired session configuration.
		// Ensure cookies, proxy protocol etc are set up if reqd.
		SessCln *http.Client
	*/

	// Browser is the headless browser to use.
	Browser *rod.Browser
)

func (arg Args) String() string {
	// 201468 c1 (group/Qvv4lz52cT)
	// 1234 (contest)
	// 100522 f1 (gym)

	var str string
	if arg.Group != "" {
		str = fmt.Sprintf("%v %v (%v/%v)", arg.Contest, arg.Problem, arg.Class, arg.Group)
	} else {
		str = fmt.Sprintf("%v %v (%v)", arg.Contest, arg.Problem, arg.Class)
	}
	return strings.Join(strings.Fields(str), " ")
}

// loginPage returns link to login page
func loginPage() string {
	return fmt.Sprintf("%v/enter", hostURL)
}

// Parse passed in specifier string to new Args struct.
// Validates parsed args and returns error if any.
//
// List of valid specifiers can be viewed at
// github.com/cp-tools/codeforces/wiki.
func Parse(str string) (Args, error) {
	var (
		rxCont  = `(?P<cont>\d+)`
		rxProb  = `(?P<prob>[A-Za-z][1-9]?)`
		rxClass = `(?P<class>contest|gym|group|problemset)`
		rxGroup = `(?P<group>\w{10})`

		valRx = []string{
			`codeforces.com\/` + rxClass + `\/` + rxCont + `$`,
			`codeforces.com\/` + rxClass + `\/` + rxCont + `\/problem\/` + rxProb + `$`,
			`codeforces.com\/` + rxClass + `\/` + rxGroup + `\/` + `contest` + `\/` + rxCont + `$`,
			`codeforces.com\/` + rxClass + `\/` + rxGroup + `\/` + `contest` + `\/` + rxCont + `\/problem\/` + rxProb + `$`,
			`codeforces.com\/` + rxClass + `\/problem\/` + rxCont + `\/` + rxProb + `$`,

			`^\s*` + rxCont + `$`,
			`^\s*` + rxCont + `\s*` + rxProb + `$`,
			`^\s*` + rxGroup + `\s*` + rxCont + `$`,
			`^\s*` + rxGroup + `\s*` + rxCont + `\s*` + rxProb + `$`,

			// for local folders parsing
			`^\s*` + rxClass + `\s*` + rxCont + `$`,
			`^\s*` + rxClass + `\s*` + rxCont + `\s*` + rxProb + `$`,
			`^\s*` + rxClass + `\s*` + rxGroup + `\s*` + rxCont + `$`,
			`^\s*` + rxClass + `\s*` + rxGroup + `\s*` + rxCont + `\s*` + rxProb + `$`,
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
			// convert to lowercase (default config)
			result["prob"] = strings.ToLower(result["prob"])
			arg := Args{
				Contest: result["cont"],
				Problem: result["prob"],
				Group:   result["group"],
			}
			arg.setContestClass()
			return arg, nil
		}
	}
	return Args{}, ErrInvalidSpecifier
}

// login tries logging into codeforces using credentials passed.
// Checks if any active session exists (in SessCln) before logging in.
// If you wish to overwrite currently logged in session, set cookies
// of SessCln to nil before logging in.
//
// If login is successful, returns user handle of now logged in session.
// Otherwise, if login fails, returns ErrInvalidCredentials as error.
//
// By default, option 'remember me' is checked, ensuring the session
// has expiry period of one month from date of last login.
func login(usr, passwd string) (string, error) {
	link := loginPage()
	page, err := Browser.PageE(link)
	if err != nil {
		return "", err
	}
	defer page.Close()

	page.WaitLoad()
	if msg := cE(page); msg != "" {
		return "", fmt.Errorf(msg)
	}

	// check if current user sesion is logged in
	if handle := findHandle(page); handle != "" {
		return handle, nil
	}

	// otherwise, login
	page.Element("#handleOrEmail").Input(usr)
	page.Element("#password").Input(passwd)
	if page.Element("#remember").Property("checked").Bool() == false {
		page.Element("#remember").Click()
	}

	wait := page.WaitRequestIdle()
	page.Element(".submit").Click()
	wait()

	if handle := findHandle(page); handle != "" {
		return handle, nil
	}
	return "", ErrInvalidCredentials
}
