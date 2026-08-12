package names

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
		result := ToCamelCase(tc.input, tc.upperFirst)
		if result != tc.expected {
			t.Errorf("Input: %s, Expected: %s, UpperFirst: %v Got: %s", tc.input, tc.expected, tc.upperFirst, result)
		}
	}
}

func TestToUnderscore(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"hello-world", "hello_world"},
		{"lorem_ipsum-dolor.sit amet!", "lorem_ipsum_dolor_sit_amet"},
		{"one!two_three4five", "one_two_three4five"},
		{"Camel-Case-Test", "camel_case_test"},
		{"ABC_def_GHI", "abc_def_ghi"},
	}
	for _, tc := range testCases {
		result := ToUnderscore(tc.input)
		if result != tc.expected {
			t.Errorf("Input: %s, Expected: %s, Got: %s", tc.input, tc.expected, result)
		}
	}
}

func TestGeneratedFilePath(t *testing.T) {
	testCases := []struct {
		path        string
		defaultDir  string
		expectedOut string
	}{
		{".gitignore", "fs", "gitignore.go"},
		{"nested/.gitignore", "fs", "nested/gitignore.go"},
		{"binaryfile.ex", "fs", "binaryfile.ex.go"},
		{"nested/file.txt", "fs", "nested/file.txt.go"},
	}
	for _, tc := range testCases {
		result := NewFileNames(tc.path, tc.defaultDir).GeneratedFilePath()
		if result != tc.expectedOut {
			t.Errorf("Path: %s, Expected: %s, Got: %s", tc.path, tc.expectedOut, result)
		}
	}
}
