package main

import (
	"bytes"
	"encoding/hex"
	"testing"
)

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
