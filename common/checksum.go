package common

import (
	"crypto/sha1"
	"fmt"
	"os"
)

func CheckSum(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha1.Sum(content)), nil
}
