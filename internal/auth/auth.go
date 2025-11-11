package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"

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

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)

	return hex.EncodeToString(key)
}

func GetAPIKey(headers http.Header) (string, error) {
	stringKey := headers.Get("Authorization")

	if stringKey == "" {
		return "", fmt.Errorf("key no found")
	}

	splitedKey := strings.Split(stringKey, " ")
	if len(splitedKey) != 2 || splitedKey[0] != "ApiKey" {
		return "", fmt.Errorf("malformed ApiKey token")
	}

	return splitedKey[1], nil
}
