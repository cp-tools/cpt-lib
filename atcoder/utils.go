package atcoder

import (
	"github.com/cp-tools/cpt-lib/util"
	"github.com/go-rod/rod/lib/proto"
)

var (
	selCSSHandle = `.navbar-right>li:last-child>a[class]`
	selCSSNotif  = `#main-container .alert`
	selCSSFooter = `footer.footer`
)

func loadPage(link string, selMatch ...string) (*page, string, error) {
	rp, err := util.NewPage(Browser, link, []proto.NetworkResourceType{
		proto.NetworkResourceTypeImage, proto.NetworkResourceTypeFont,
		proto.NetworkResourceTypeStylesheet, proto.NetworkResourceTypeMedia,
	})
	if err != nil {
		return nil, "", err
	}

	p := &page{rp}

	rc := p.Race()
	rc.Element(selCSSNotif)
	for _, sel := range selMatch {
		rc.Element(sel)
	}

	elm := rc.MustDo()
	if elm.MustMatches(selCSSNotif) {
		return p, util.StrClean(elm.MustText()), nil
	}

	return p, "", nil
}
