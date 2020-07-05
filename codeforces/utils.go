package codeforces

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func parseResp(resp *http.Response) ([]byte, string) {
	var msg string
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	msgRgx := `Codeforces\.showMessage\("(.+)"\);\s+Codeforces\.reformatTimes\(\);`
	re := regexp.MustCompile(msgRgx)
	tmp := re.FindStringSubmatch(string(body))
	if tmp != nil {
		msg = clean(tmp[1])
	}
	return body, msg
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
	// replace non-ascii with space
	re = regexp.MustCompile(`[\x00\x08\x0B\x0C\x0E-\x1F]+`)
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
func findHandle(body []byte) string {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	val := doc.Find("#header").Find("a[href^=\"/profile/\"]").Text()
	return val
}

// findCsrf extracts Csrf from REQUEST body
func findCsrf(body []byte) string {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	val, _ := doc.Find(".csrf-token").Attr("data-csrf")
	return val
}

// findPagination returns number of pages of table
// returns (1 if no pagination found)
func findPagination(body []byte) int {
	// parse html body to find number of pages (in pagination)
	// return's default value of 1 if no pagination found
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	val := getText(doc.Find(".page-index").Last(), "a")
	num, err := strconv.Atoi(val)
	if err != nil {
		return 1
	}
	return num
}

// if the time string is invalid, returns time corresponding to
// the start of time => (1 Jan 1970 00:00)
func parseTime(str string) time.Time {
	// date-time format on codeforces
	const ruTime = "02.01.2006 15:04 Z07:00"
	const enTime = "Jan/02/2006 15:04 Z07:00"

	raw := fmt.Sprintf("%v +03:00", str)
	tm, err := time.Parse(enTime, raw)
	if err != nil {
		tm, err = time.Parse(ruTime, raw)
		if err != nil {
			// set to the beginning of time
			tm = time.Unix(0, 0)
		}
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
