package utils

import (
	"strings"
)

func LowercaseFirstLetter(s string) string {
	if len(s) > 0 {
		return strings.ToLower(s[0:1]) + s[1:]
	}
	return s
}

func UppercaseFirstLetter(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(s[0:1]) + s[1:]
	}
	return s
}
