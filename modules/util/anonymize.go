package util

import (
	"strings"
)

// ObfuscateEmail hides away most of an email, prior to @ and returns the result
//
// E.g:
//
//	ObfuscateEmail("helloYou@gmail.com") == "h...@gmail.com"
func ObfuscateEmail(email string) string {
	if !strings.Contains(email, "@") {
		return ObfuscateString(email, false)
	}

	split := strings.SplitN(email, "@", 2)
	return split[0][:1] + "...@" + split[1]
}

func ObfuscateString(value string, isPhone bool) string {
	strLen := len(value)
	cutoff := 1
	if isPhone {
		cutoff = 3
	}
	return value[:cutoff] + strings.Repeat(".", strLen-cutoff)
}
