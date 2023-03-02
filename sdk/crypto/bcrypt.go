package crypto

import "golang.org/x/crypto/bcrypt"

func HashValue(value string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHash(hash, value string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
}
