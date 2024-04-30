package names

import (
	"strings"
	"unicode"
)

func isNameSeparator(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsNumber(r)
}

// ToCamelCase convert string to CamelCase name
func ToCamelCase(s string, upperFirst bool) string {
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

// ToUnderscore convert string to underscore name
func ToUnderscore(s string) string {
	var name strings.Builder
	var isNewWord = false
	for _, char := range s {
		if isNameSeparator(char) {
			isNewWord = true
		} else {
			if isNewWord {
				isNewWord = false
				name.WriteRune('_')
			}
			name.WriteRune(unicode.ToLower(char))
		}
	}
	return name.String()
}
