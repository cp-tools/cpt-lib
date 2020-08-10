package atcoder

import (
	"regexp"
	"strings"

	"github.com/go-rod/rod"
)

func cE(page *rod.Page) string {
	elm := page.Elements(".alert-danger").First()
	if elm == nil {
		return ""
	}
	return clean(elm.Text())
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

func findHandle(page *rod.Page) string {
	elm := page.Elements(".navbar-right > li")
	if len(elm) != 2 {
		return ""
	}
	handle := elm.Last().Elements("a").First().Text()
	return clean(handle)
}
