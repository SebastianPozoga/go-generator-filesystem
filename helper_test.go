package main

import (
	"testing"
)

func TestToCamelCase(t *testing.T) {
	testCases := []struct {
		input      string
		expected   string
		upperFirst bool
	}{
		{"hello-world", "helloWorld", false},
		{"lorem_ipsum-dolor.sit amet!", "loremIpsumDolorSitAmet", false},
		{"one!two_three4five", "oneTwoThree4five", false},
		{"Camel-Case-Test", "camelCaseTest", false},
		{"ABC_def_GHI", "abcDefGhi", false},
		{"hello-world", "HelloWorld", true},
		{"lorem_ipsum-dolor.sit amet!", "LoremIpsumDolorSitAmet", true},
		{"one!two_three4five", "OneTwoThree4five", true},
		{"Camel-Case-Test", "CamelCaseTest", true},
		{"ABC_def_GHI", "AbcDefGhi", true},
	}
	for _, tc := range testCases {
		result := toCamelCase(tc.input, tc.upperFirst)
		if result != tc.expected {
			t.Errorf("Input: %s, Expected: %s, UpperFirst: %v Got: %s", tc.input, tc.expected, tc.upperFirst, result)
		}
	}
}
