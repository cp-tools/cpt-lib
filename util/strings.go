package util

import (
	"math/rand"
	"regexp"
)

// StrClean replaces all unicode space
// characters in the string with ascii space.
func StrClean(str string) string {
	re := regexp.MustCompile(`\p{Z}`)
	return re.ReplaceAllString(str, " ")
}

// StrRandom generates a random string of length n.
// The returned string is strictly alpha-numeric.
// Attrib: https://stackoverflow.com/a/31832326/9606036.
func StrRandom(n int) string {
	const charBytes = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = charBytes[rand.Intn(len(charBytes))]
	}

	return string(b)
}
