package utils

import "golang.org/x/crypto/bcrypt"

func NormalizedPassword(p string) []byte {
	return []byte(p)
}

func GeneratePassword(p string) string {
	password := NormalizedPassword(p)

	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)

	if err != nil {
		return err.Error()
	}

	return string(hash)
}

func ComparePassword(password string, hash string) bool {
	byteHash := NormalizedPassword(hash)
	bytePwd := NormalizedPassword(password)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	return err == nil
}
