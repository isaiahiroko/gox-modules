package utils

import (
	"crypto/sha512"
	"encoding/hex"
)

func Checksum(data string) string {
	hash := sha512.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
