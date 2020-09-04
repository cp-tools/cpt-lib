package codeforces

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
)

var (
	selCSSNotif   = `.jGrowl-notification .message`
	selCSSHandle  = `#header a[href^="/profile/"]`
	selCSSCurrTab = `.second-level-menu-list .current`
	selCSSFooter  = `#footer`
	selCSSError   = `.error`
)

func loadPage(link string) (*rod.Page, string, error) {
	page, err := Browser.Page(link)
	if err != nil {
		return nil, "", err
	}

	// footer is loaded last ig? I'm not sure
	elm := page.MustElement(selCSSNotif, selCSSFooter)
	if elm.MustMatches(selCSSNotif) {
		return page, clean(elm.MustText()), nil
	}
	return page, "", nil
}

func clean(str string) string {
	// remove trailiing/leading spaces
	str = strings.ReplaceAll(str, "<br/>", "\n")
	str = strings.TrimSpace(str)
	// remove extra whitespace after \n
	re := regexp.MustCompile(`\n\s+`)
	str = re.ReplaceAllString(str, "\n")
	// remove extra whitespaces
	re = regexp.MustCompile(` +`)
	str = re.ReplaceAllString(str, " ")
	// replace any space character space
	re = regexp.MustCompile(`\p{Z}`)
	return re.ReplaceAllString(str, " ")
}

// getText extracts text from particular html data
func getText(sel *goquery.Selection, query string) string {
	str := sel.Find(query).Text()
	return clean(str)
}

// getAttr extracts attribute valur of particular html data
func getAttr(sel *goquery.Selection, query, attr string) string {
	str := sel.Find(query).AttrOr(attr, "")
	return clean(str)
}

// if the time string is invalid, returns time corresponding to
// the start of time => (1 Jan 1970 00:00)
func parseTime(link string) time.Time {
	re := regexp.MustCompile(`([A-Za-z]+)\/(\d+)\/(\d+) (\d+):(\d+)`)
	pst := re.FindAllStringSubmatch(link, -1)
	if pst == nil || len(pst[0]) < 6 {
		return time.Unix(0, 0).UTC()
	}

	// set values
	pMonth, pDay, pYear := pst[0][1], pst[0][2], pst[0][3]
	pHour, pMinute := pst[0][4], pst[0][5]
	val := fmt.Sprintf("%v/%v/%v %v:%v",
		pMonth, pDay, pYear, pHour, pMinute)

	// only if UTC... is present
	re = regexp.MustCompile(`UTC(\+|-)(\d+).(\d+)`)
	pst = re.FindAllStringSubmatch(link, -1)
	if pst == nil || len(pst[0]) < 4 {
		val = fmt.Sprintf("%v +00:00", val)
	} else {
		pOffset, pMajor, pMinor := pst[0][1], pst[0][2], pst[0][3]
		pMajor = fmt.Sprintf("0%v", pMajor)[:2]
		if pMinor == "5" {
			pMinor = "30"
		}
		pMinor = fmt.Sprintf("%v0", pMinor)[:2]

		val = fmt.Sprintf("%v %v%v:%v", val, pOffset, pMajor, pMinor)
	}

	tm, err := time.Parse("Jan/2/2006 15:04 Z07:00", val)
	if err != nil {
		tm = time.Unix(0, 0)
	}
	return tm.UTC()
}

// genRandomString generates a random string of length n.
// Code copied from https://stackoverflow.com/a/9606036.
func genRandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (arg *Args) setContestClass() {
	val, err := strconv.Atoi(arg.Contest)
	if len(arg.Group) == 10 {
		arg.Class = ClassGroup
	} else if err == nil {
		if val <= 100000 {
			arg.Class = ClassContest
		} else {
			arg.Class = ClassGym
		}
	}
}
