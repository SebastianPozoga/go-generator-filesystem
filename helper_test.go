package main

import (
	"bytes"
	"encoding/hex"
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

func TestCalculateChecksum(t *testing.T) {
	var (
		testData                   = []byte("This is a test file.")
		checksum, expectedChecksum []byte
		err                        error
	)
	if expectedChecksum, err = hex.DecodeString("f29bc64a9d3732b4b9035125fdb3285f5b6455778edca72414671e0ca3b2e0de"); err != nil {
		panic(err)
	}
	if checksum, err = calculateChecksum(testData); err != nil {
		t.Fatal("Error calculating checksum:", err)
	}
	if !bytes.Equal(checksum, expectedChecksum) {
		t.Errorf("Expected checksum: %v, got: %v", expectedChecksum, checksum)
	}
}
