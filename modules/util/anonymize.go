package util

import (
	"strings"
)

func ObfuscateEmail(email string) string {
	splitString := strings.Split(email, "@")
	return splitString[0][:1] + "...@" + splitString[1]
}

func ObfuscateString(value string, isPhone bool) string {
	strLen := len(value)
	cutoff := 1
	if isPhone {
		cutoff = 3
	}
	return value[:cutoff] + strings.Repeat(".", strLen-cutoff)
}
