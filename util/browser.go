package util

import (
	"os"
	"path/filepath"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// NewBrowser initiates the automated browser to use.
func NewBrowser(headless bool, userDataDir, bin string) (*rod.Browser, error) {
	// Launch browser.
	launchBrowser := func(controlURL string) (*rod.Browser, error) {
		b := rod.New().ControlURL(controlURL)
		if err := b.Connect(); err != nil {
			return nil, err
		}
		return b, nil
	}

	// Store data in cache (to reduce time).
	cacheDir, _ := os.UserCacheDir()
	cacheUserDataDir := filepath.Join(cacheDir, "cp-tools", "cpt-lib", bin)

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
	if file, err := os.Stat(userDataDir); err == nil && file.IsDir() {
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
