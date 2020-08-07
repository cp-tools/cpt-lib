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

func cE(page *rod.Page) string {
	var msg string
	body := page.Element("html").HTML()

	msgRgx := `Codeforces\.showMessage\("(.+)"\);\s+Codeforces\.reformatTimes\(\);`
	re := regexp.MustCompile(msgRgx)
	tmp := re.FindStringSubmatch(body)
	if tmp != nil {
		msg = clean(tmp[1])
	}
	return msg
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

// findHandle scrapes handle from REQUEST body
func findHandle(page *rod.Page) string {
	elm := page.Elements("#header a[href^=\"/profile/\"]")
	if len(elm) == 0 {
		return ""
	}

	return elm.First().Text()
}

// findPagination returns number of pages of table
// returns (1 if no pagination found)
func findPagination(page *rod.Page) int {
	val := page.Elements(".page-index a").Last()
	if val == nil {
		// no pagination found
		return 1
	}
	num, _ := strconv.Atoi(val.Text())
	return num
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
