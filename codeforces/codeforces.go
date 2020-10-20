package codeforces

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type (
	// Args holds specifier details parsed by
	// Parse() function.
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
	errInvalidCredentials = fmt.Errorf("invalid login credentials")
)

var (
	hostURL = "https://codeforces.com"

	// Browser is the headless browser to use.
	Browser *rod.Browser
)

// Start initiates the headless browser to use.
func Start(headless bool, userDataDir, bin string, flags ...[]string) error {
	// Add flags to launcher.
	addFlags := func(l *launcher.Launcher) {
		for _, flag := range flags {
			l.Set(flag[0], flag[1:]...)
		}
	}
	// Launch browser.
	launchBrowser := func(controlURL string) (*rod.Browser, error) {
		b := rod.New().ControlURL(controlURL)
		if err := b.Connect(); err != nil {
			return nil, err
		}
		return b, nil
	}

	// Initiate browser to extract cookies from.
	cookiesl := launcher.NewUserMode().
		UserDataDir(userDataDir).
		Headless(true).
		Bin(bin)
	addFlags(cookiesl)

	cookiesControlURL, err := cookiesl.Launch()
	if err != nil {
		return err
	}

	cookiesBrowser, err := launchBrowser(cookiesControlURL)
	if err != nil {
		return err
	}
	defer cookiesBrowser.Close()

	// Initiate the browser to use.
	l := launcher.New().
		UserDataDir(filepath.Join(userDataDir, "cpt-lib")).
		Headless(headless).
		Bin(bin)
	addFlags(l)

	controlURL, err := l.Launch()
	if err != nil {
		return err
	}

	Browser, err = launchBrowser(controlURL)
	if err != nil {
		return err
	}

	// Copy cookies of user.
	Browser.MustSetCookies(cookiesBrowser.MustGetCookies())

	return nil
}

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
			`codeforces.com\/problemset\/problem\/` + rxCont + `\/` + rxProb + `$`,

			`^\s*` + rxClass + `$`,
			`^\s*` + rxGroup + `$`,

			`^\s*` + rxCont + `$`,
			`^\s*` + rxCont + `\s*` + rxProb + `$`,
			`^\s*` + rxGroup + `\s*` + rxCont + `$`,
			`^\s*` + rxGroup + `\s*` + rxCont + `\s*` + rxProb + `$`,
		}
	)

	str = strings.TrimSpace(str)
	if str == "" {
		return Args{}, nil
	}

	for _, rgx := range valRx {
		re := regexp.MustCompile(rgx)
		if re.MatchString(str) {
			// attrib : stackoverflow.com/a/9606036
			match := re.FindStringSubmatch(str)
			result := map[string]string{}
			for i, name := range re.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}
			// convert to lowercase (default config)
			result["prob"] = strings.ToLower(result["prob"])
			arg := Args{
				Contest: result["cont"],
				Problem: result["prob"],
				Class:   result["class"],
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
	page, msg := loadPage(link, selCSSFooter)
	defer page.Close()

	if msg != "" {
		// there shouldn't be any notification
		return "", fmt.Errorf(msg)
	}

	// check if current user sesion is logged in
	if elm := page.MustElements(selCSSHandle).First(); elm != nil {
		return clean(elm.MustText()), nil
	}
	// otherwise, login

	// check if username/password are valid
	if usr == "" || passwd == "" {
		return "", errInvalidCredentials
	}

	page.MustElement("#handleOrEmail").Input(usr)
	page.MustElement("#password").Input(passwd)
	if page.MustElement("#remember").MustProperty("checked").Bool() == false {
		page.MustElement("#remember").MustClick()
	}
	page.MustElement(".submit").MustClick()

	page.MustElement(selCSSError, selCSSHandle)
	if elm := page.MustElements(selCSSHandle); !elm.Empty() {
		return clean(elm.First().MustText()), nil
	}

	return "", errInvalidCredentials
}

func logout() error {
	page, msg := loadPage(hostURL, selCSSFooter)
	defer page.Close()

	if msg != "" {
		return fmt.Errorf(msg)
	}

	if page.MustHasR("a", "Logout") {
		page.MustElementR("a", "Logout").MustClick()
		// page gives a notification on logout
		page.Element(selCSSNotif)
	}
	return nil
}
