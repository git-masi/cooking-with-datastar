package internal

import (
	"strings"
	"unicode"
)

const (
	HYPHEN     = '-'
	UNDERSCORE = '_'
)

func ToStartCase(str string) string {
	start := 0

	for i, v := range str {
		if unicode.IsLetter(v) {
			start = i
			break
		}
	}

	if start == len(str)-1 {
		// No letters, return early
		return str
	}

	prev := 0
	parts := []string{}

	for i, v := range str {
		if i < start {
			continue
		}

		if i == 0 && unicode.IsUpper(v) {
			// If the first rune in the string is an uppercase letter continue
			continue
		}

		if v == HYPHEN || v == UNDERSCORE || unicode.IsSpace(v) {
			parts = append(parts, str[prev:i])
			prev = i + 1
			continue
		}

		if unicode.IsUpper(v) {
			parts = append(parts, str[prev:i])
			prev = i

		}
	}

	if prev < len(str) {
		parts = append(parts, str[prev:])
	}

	for i, v := range parts {
		if i == 0 {
			parts[i] = strings.Join([]string{strings.ToUpper(v[0:1]), v[1:]}, "")
		} else {
			parts[i] = strings.Join([]string{strings.ToLower(v[0:1]), v[1:]}, "")
		}
	}

	return strings.Join(parts, " ")
}
