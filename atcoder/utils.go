package atcoder

import (
	"regexp"
	"strings"

	"github.com/go-rod/rod"
)

var (
	selCSSHandle = `.navbar-right>li:last-child a`
	selCSSFooter = `footer.footer`
	selCSSNotif  = `.alert`
)

func loadPage(link string) (*rod.Page, string, error) {
	page, err := Browser.PageE(link)
	if err != nil {
		return nil, "", err
	}

	// footer is loaded last ig? not sure
	elm := page.Element(selCSSNotif, selCSSFooter)
	if elm.Matches(selCSSNotif) {
		return page, clean(elm.Text()), nil
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
