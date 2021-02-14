package codeforces

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/cp-tools/cpt-lib/v2/util"

	"github.com/go-rod/rod"
)

type (
	// Args holds specifier details parsed by
	// Parse() function. All methods use this
	// at the core.
	Args struct {
		Contest string
		Problem string
		Class   string
		Group   string
	}

	page struct {
		*rod.Page
	}
)

// Class type of contest.
const (
	ClassContest = "contest"
	ClassGroup   = "group"
	ClassGym     = "gym"
)

// Errors returned by library.
var (
	ErrInvalidSpecifier = fmt.Errorf("invalid specifier data")
)

var (
	hostURL = "https://codeforces.com"

	// Browser is the headless browser to use.
	Browser *rod.Browser
)

func (arg Args) String() (str string) {
	if arg == (Args{}) {
		return ""
	}

	switch arg.Class {
	case ClassGroup:
		str = fmt.Sprintf("%v %v (%v/%v)", arg.Contest, arg.Problem, arg.Class, arg.Group)

	case ClassContest, ClassGym:
		str = fmt.Sprintf("%v %v (%v)", arg.Contest, arg.Problem, arg.Class)
	}

	return strings.Join(strings.Fields(str), " ")
}

// Start initiates the automated browser to use.
func Start(headless bool, userDataDir, bin string) error {
	cacheDir, _ := os.UserCacheDir()
	cacheDir = filepath.Join(cacheDir, "cp-tools", "cpt-lib")

	return StartWithCacheDir(headless, userDataDir, bin, cacheDir)
}

// StartWithCacheDir is the same as Start, only allows to set cacheDir to use.
func StartWithCacheDir(headless bool, userDataDir, bin, cacheDir string) error {
	bs, err := util.NewBrowser(headless, userDataDir, bin, cacheDir)
	Browser = bs

	return err
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

			// Set contest class.
			switch {
			case arg.Class != "":
				break
			case len(arg.Group) == 10:
				arg.Class = ClassGroup
			default:
				switch val, _ := strconv.Atoi(arg.Contest); {
				case val <= 1e5:
					arg.Class = ClassContest
				default:
					arg.Class = ClassGym
				}
			}

			return arg, nil
		}
	}
	return Args{}, ErrInvalidSpecifier
}
