package usaco

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// setup headless browser to use
	_, mode := os.LookupEnv("LOCAL_MODE")

	Start(!mode, "../user-data-dir", "google-chrome",
		[]string{"disable-extensions"})

	// setup login access to use
	usr := os.Getenv("USACO_USERNAME")
	passwd := os.Getenv("USACO_PASSWORD")
	_, err := login(usr, passwd)
	if err != nil {
		fmt.Println("Login failed:", err)
		Browser.Close()
		os.Exit(1)
	}

	exitCode := m.Run()

	// logout current user
	if err := logout(); err != nil {
		fmt.Println("Logout failed:", err)
		Browser.Close()
		os.Exit(1)
	}

	Browser.Close()
	os.Exit(exitCode)
}
