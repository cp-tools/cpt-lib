package usaco

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

var (
	selCSSNavbar  = `.navbar`
	selCSSFormErr = `.form_error`
	selCSSHandle  = `div:nth-child(5)>div:nth-child(1)>strong`
)

func loadPage(link string) (*rod.Page, error) {
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

	page.Element(selCSSNavbar)
	return page, nil
}

func processHTML(page *rod.Page) *goquery.Document {
	doc, _ := goquery.NewDocumentFromReader(
		strings.NewReader(page.MustElement("html").MustHTML()))
	return doc
}

func clean(str string) string {
	// remove trailing/leading spaces
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
