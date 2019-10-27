package yagolib

import (
	"strings"
	"unicode"
)

// IsFirstRuneUpper returns 'true' if the 'str' begins with a capital letter.
func IsFirstRuneUpper(str string) bool {
	for _, r := range str {
		return r == unicode.ToUpper(r)
	}
	return false
}

// RemoveCharacters removes all chars contained in 'chars' from input string 's'.
func RemoveCharacters(s, chars string) string {
	filter := func(r rune) rune {
		if strings.IndexRune(chars, r) < 0 {
			return r
		}
		return -1
	}
	return strings.Map(filter, s)
}

// GetBaseOfIntString returns the base of integer value contained in input string.
// The function doesn't check the possibility of conversion string to integer.
// The current version recognizes hex and binary literals only.
func GetBaseOfIntString(intStr string) int {
	intStr = strings.ToLower(intStr)
	if strings.HasPrefix(intStr, "0x") || strings.HasSuffix(intStr, "h") {
		return 16
	} else if strings.HasPrefix(intStr, "0b") || strings.HasSuffix(intStr, "b") {
		return 2
	}
	return 10
}
