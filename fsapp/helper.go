package fsapp

import (
	"crypto/sha256"
)

func calculateChecksum(bytes []byte) (result []byte, err error) {
	hash := sha256.New()
	if _, err = hash.Write(bytes); err != nil {
		return
	}
	return hash.Sum(nil), nil
}
