package utils

import (
	"encoding/base64"
)

func EncodeBase64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func DecodeBase64(data string) string {
	result, _ := base64.StdEncoding.DecodeString(data)

	return string(result)
}
