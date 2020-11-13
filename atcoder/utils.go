package atcoder

import (
	"github.com/cp-tools/cpt-lib/util"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

var (
	selCSSHandle = `.navbar-right>li:last-child a`
	selCSSNotif  = `#main-container .alert`
	selCSSFooter = `footer.footer`
)

func loadPage(link string, selMatch ...string) (*rod.Page, string, error) {
	page, err := util.NewPage(Browser, link, []proto.NetworkResourceType{
		proto.NetworkResourceTypeImage, proto.NetworkResourceTypeFont,
		proto.NetworkResourceTypeStylesheet, proto.NetworkResourceTypeMedia,
	})
	if err != nil {
		return nil, "", err
	}

	selMatch = append([]string{selCSSNotif}, selMatch...)
	elm := page.MustElement(selMatch...)

	if elm.MustMatches(selCSSNotif) {
		return page, elm.MustText(), nil
	}

	return page, "", nil
}
