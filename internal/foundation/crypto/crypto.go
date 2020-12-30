package crypto

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func HashValue(p string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrapf(err, "HashValue -> bcrypt.GenerateFromPassword(***, %v)", bcrypt.DefaultCost)
	}
	return string(hash), nil
}

func CompareHashes(plain string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		return false
	}
	return true
}
