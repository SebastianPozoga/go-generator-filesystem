package main

import (
	"strings"
	"unicode"
)

func isNameSeparator(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsNumber(r)
}

// toCamelCase convert string to CamelCase name
func toCamelCase(s string, upperFirst bool) string {
	var name strings.Builder
	var isNewWord = upperFirst
	for _, char := range s {
		if isNameSeparator(char) {
			isNewWord = true
		} else {
			if isNewWord {
				isNewWord = false
				name.WriteRune(unicode.ToUpper(char))
				continue
			}
			name.WriteRune(unicode.ToLower(char))
		}
	}
	return name.String()
}

// toUnderscore convert string to underscore name
func toUnderscore(s string) string {
	var name strings.Builder
	var isNewWord = false
	for _, char := range s {
		if isNameSeparator(char) {
			isNewWord = true
		} else {
			if isNewWord {
				isNewWord = false
				name.WriteRune('_')
				continue
			}
			name.WriteRune(unicode.ToLower(char))
		}
	}
	return name.String()
}
