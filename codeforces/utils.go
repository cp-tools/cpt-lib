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
	"github.com/go-rod/rod/lib/proto"
)

var (
	selCSSNotif  = `.jGrowl-notification .message`
	selCSSHandle = `#header a[href^="/profile/"]`
	selCSSFooter = `#footer`
	selCSSError  = `.error`
)

func loadPage(link string, selMatch ...string) (*rod.Page, string) {
	page := Browser.MustPage(link)
	// Disable CSS and Img in webpage.
	router := page.HijackRequests()
	router.MustAdd("*", func(h *rod.Hijack) {
		if h.Request.Type() == proto.NetworkResourceTypeImage ||
			h.Request.Type() == proto.NetworkResourceTypeFont ||
			h.Request.Type() == proto.NetworkResourceTypeStylesheet ||
			h.Request.Type() == proto.NetworkResourceTypeMedia {
			h.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}
		h.ContinueRequest(&proto.FetchContinueRequest{})
	})
	go router.Run()

	selMatch = append([]string{selCSSNotif}, selMatch...)
	elm := page.MustElement(selMatch...)

	if page.MustInfo().URL != link {
		page.WaitLoad()
		if page.MustHas(selCSSNotif) {
			// There was a redirect (with an error message).
			elm = page.MustElement(selCSSNotif)
		}
	}

	if elm.MustMatches(selCSSNotif) {
		return page, clean(elm.MustText())
	}
	return page, ""
}

func processHTML(page *rod.Page) *goquery.Document {
	doc, _ := goquery.NewDocumentFromReader(
		strings.NewReader(page.MustElement("html").MustHTML()))
	return doc
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
	// Prepare for data extraction (strip all extra whitespace).
	link = strings.ReplaceAll(clean(link), "\n", " ")

	// Follows english locale format: Mon/dd/yyyy hh:mm +MM:mm
	re := regexp.MustCompile(`([A-Za-z]{3})\/(\d{2})\/(\d{4}) (\d+):(\d+)`)
	pst := re.FindAllStringSubmatch(link, -1)
	if pst == nil || len(pst[0]) < 6 {
		// Try the russian locale format: dd.mm.yyyy hh:mm +MM:mm
		re = regexp.MustCompile(`(\d{2})\.(\d{2})\.(\d{4}) (\d+):(\d+)`)
		pst = re.FindAllStringSubmatch(link, -1)
		if pst == nil || len(pst[0]) < 6 {
			// Formats didn't match. Mostly invalid.
			return time.Unix(0, 0).UTC()
		}
		// Convert month int to short string name.
		mm, _ := strconv.Atoi(pst[0][2])
		mon := time.Month(mm).String()[:3]
		// Rearrange pst values to match english locale.
		pst[0][2], pst[0][1] = pst[0][1], mon
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
		pMajor = fmt.Sprintf("0%v", pMajor)
		pMajor = pMajor[len(pMajor)-2:]

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
	if arg.Class != "" {
		return
	}

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
