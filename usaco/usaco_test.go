package usaco

import (
	"os"
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func TestMain(m *testing.M) {
	// setup headless browser to use
	_, mode := os.LookupEnv("LOCAL_MODE")

	l := launcher.New().UserDataDir("../user-data-dir").
		Set("blink-settings", "imagesEnabled=false")
	if mode {
		l.Headless(false)
	}
	Browser = rod.New().ControlURL(l.Launch()).Connect()

	if !mode {
		// setup login access to use
		usr := os.Getenv("USACO_USERNAME")
		passwd := os.Getenv("USACO_PASSWORD")
		login(usr, passwd)
	}

	exitCode := m.Run()

	// logout current user
	if !mode {
		logout()
	}
	// close the browser instance
	Browser.Close()

	os.Exit(exitCode)
}
