package auth

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func HashPassword(pass string) (string, error) {
	hash, err := argon2id.CreateHash(pass, argon2id.DefaultParams)

	if err != nil {
		log.Fatalf(`[error] argon2id.CreateHash(pass, argon2id.DefaultParams): %v`, err)
		return "", err
	}

	return hash, nil
}

func CheckPasswordHash(pass string, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(pass, hash)
	if err != nil {
		log.Fatalf(`[error] argon2id.ComparePasswordAndHash(pass, hash): %v`, err)
		return false, err
	}

	return match, nil
}
