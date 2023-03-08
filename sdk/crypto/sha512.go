package crypto

import (
	"crypto/sha512"
	"encoding/hex"
)

func HashValueSha512(value string) (string, error) {
	hash := sha512.New()
	if _, err := hash.Write([]byte(value)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func CompareHashSha512(hash, value string) error {
	hashCompare, err := HashValueSha512(value)
	if err != nil {
		return err
	}
	if hash != hashCompare {
		return err
	}
	return nil
}
