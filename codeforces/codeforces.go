package codeforces

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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
	// SessCln should be set to desired session configuration.
	// Ensure cookies, proxy protocol etc are set up if reqd.
	SessCln *http.Client
)

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
		rxClass = `(?P<class>contest|gym|group)`
		rxGroup = `(?P<group>\w{10})`

		valRx = []string{
			`codeforces.com\/` + rxClass + `\/` + rxCont + `$`,
			`codeforces.com\/` + rxClass + `\/` + rxCont + `\/problem\/` + rxProb + `$`,
			`codeforces.com\/` + rxClass + `\/` + rxGroup + `\/` + `contest` + `\/` + rxCont + `$`,
			`codeforces.com\/` + rxClass + `\/` + rxGroup + `\/` + `contest` + `\/` + rxCont + `\/problem\/` + rxProb + `$`,

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

// Login tries logging into codeforces using credentials passed.
// Checks if any active session exists (in SessCln) before logging in.
// If you wish to overwrite currently logged in session, set cookies
// of SessCln to nil before logging in.
//
// If login is successful, returns user handle of now logged in session.
// Otherwise, if login fails, returns ErrInvalidCredentials as error.
//
// By default, option 'remember me' is checked, ensuring the session
// has expiry period of one month from date of last login.
func Login(usr, passwd string) (string, error) {
	link := loginPage()
	resp, err := SessCln.Get(link)
	if err != nil {
		return "", err
	}
	body, msg := parseResp(resp)
	if len(msg) != 0 {
		return "", err
	}

	// check if current user sesion is logged in
	if handle := findHandle(body); len(handle) != 0 {
		return handle, nil
	}

	// hidden form data
	csrf := findCsrf(body)
	ftaa := genRandomString(18)
	bfaa := genRandomString(32)

	resp, err = SessCln.PostForm(link, url.Values{
		"csrf_token":    {csrf},
		"action":        {"enter"},
		"ftaa":          {ftaa},
		"bfaa":          {bfaa},
		"handleOrEmail": {usr},
		"password":      {passwd},
		"_tta":          {"176"},
		"remember":      {"on"},
	})
	if err != nil {
		return "", err
	}

	// the only message possible is Welcome, handle!
	body, _ = parseResp(resp)
	handle := findHandle(body)
	if len(handle) == 0 {
		// login failed
		return "", ErrInvalidCredentials
	}
	return handle, nil
}
