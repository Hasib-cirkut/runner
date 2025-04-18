package utils

import "encoding/base64"

func Base64Decode(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		return "", err
	}

	return string(decoded), nil
}
