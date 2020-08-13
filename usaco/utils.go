package usaco

import (
	"regexp"
	"strings"

	"github.com/go-rod/rod"
)

var (
	selCSSNavbar  = `.navbar`
	selCSSFormErr = `.form_error`
	selCSSHandle  = `div:nth-child(5)>div:nth-child(1)>strong`
)

func loadPage(link string) (*rod.Page, error) {
	page, err := Browser.PageE(link)
	if err != nil {
		return nil, err
	}

	page.Element(selCSSNavbar)
	return page, nil
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
