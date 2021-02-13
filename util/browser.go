package util

import (
	"os"
	"path/filepath"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// NewBrowser initiates the automated browser to use.
func NewBrowser(headless bool, userDataDir, bin, cacheDir string) (*rod.Browser, error) {
	// Launch browser.
	launchBrowser := func(controlURL string) (*rod.Browser, error) {
		b := rod.New().ControlURL(controlURL)
		if err := b.Connect(); err != nil {
			return nil, err
		}
		return b, nil
	}

	// Store data in cache (to reduce time).
	cacheUserDataDir := filepath.Join(cacheDir, bin)

	// Initiate the browser to use.
	l := launcher.New().
		UserDataDir(cacheUserDataDir).
		Headless(headless).
		Bin(bin)

	controlURL, err := l.Launch()
	if err != nil {
		return nil, err
	}

	Browser, err := launchBrowser(controlURL)
	if err != nil {
		return nil, err
	}

	// Load temporary browser to extract cookies only if path exists.
	if file, err := os.Stat(userDataDir); userDataDir != "" && err == nil && file.IsDir() {
		// Initiate browser to extract cookies from.
		cookiesl := launcher.NewUserMode().
			UserDataDir(userDataDir).
			Headless(true).
			Bin(bin)

		cookiesControlURL, err := cookiesl.Launch()
		if err != nil {
			return nil, err
		}

		cookiesBrowser, err := launchBrowser(cookiesControlURL)
		if err != nil {
			return nil, err
		}
		defer cookiesBrowser.Close()
		// Copy cookies of user.
		Browser.MustSetCookies(cookiesBrowser.MustGetCookies())
	}

	return Browser, nil
}

// NewPage loads the given link in a new browser tab.
func NewPage(browser *rod.Browser, link string, block []proto.NetworkResourceType) (*rod.Page, error) {
	page, err := browser.Page(proto.TargetCreateTarget{URL: link})
	if err != nil {
		return nil, err
	}

	router := page.HijackRequests()
	router.MustAdd("*", func(h *rod.Hijack) {
		for _, b := range block {
			if h.Request.Type() == b {
				h.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
				return
			}
		}
		h.ContinueRequest(&proto.FetchContinueRequest{})
	})
	go router.Run()

	return page, nil
}
