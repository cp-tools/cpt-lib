package atcoder

import (
	"github.com/cp-tools/cpt-lib/util"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

var (
	selCSSHandle = `.navbar-right>li:last-child>a[class]`
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

	rc := page.Race()
	rc.Element(selCSSNotif)
	for _, sel := range selMatch {
		rc.Element(sel)
	}

	elm := rc.MustDo()
	if elm.MustMatches(selCSSNotif) {
		return page, util.StrClean(elm.MustText()), nil
	}

	return page, "", nil
}
