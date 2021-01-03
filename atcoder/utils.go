package atcoder

import (
	"fmt"

	"github.com/cp-tools/cpt-lib/util"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func loadPage(link string) (*page, error) {
	// Blocking these files results in faster page loading.
	resourcesToBlock := []proto.NetworkResourceType{
		proto.NetworkResourceTypeFont,
		proto.NetworkResourceTypeMedia,
		proto.NetworkResourceTypeImage,
		proto.NetworkResourceTypeStylesheet,
	}

	p, err := util.NewPage(Browser, link, resourcesToBlock)
	return &page{p}, err
}

func handleErrMsg(e *rod.Element) error {
	// There should be no notification.
	return fmt.Errorf(e.MustText())
}
