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

	page struct {
		*rod.Page
	}
)

// Errors returned by library.
var (
	ErrInvalidSpecifier = fmt.Errorf("invalid specifier data")
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
