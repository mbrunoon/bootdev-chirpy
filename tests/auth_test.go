package tests

import (
	"testing"

	"github.com/mbrunoon/bootdev-chirpy/internal/auth"
)

func TestPasswordHash(t *testing.T) {
	pass := "test"
	hashedPassword, err := auth.HashPassword(pass)

	if err != nil {
		t.Errorf(`[TestPasswordHash] auth.HashPassword(pass): %v`, err)
	}

	if hashedPassword == "" {
		t.Errorf(`[TestPasswordHash]: hashedPassword is empty`)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	pass := "test"
	hashedPassword, _ := auth.HashPassword(pass)

	if hashedPassword == "" {
		t.Errorf(`[TestCheckPasswordHash] result expected hashedPassword not empty, it was: %v`, hashedPassword)
	}

	matchFalse, _ := auth.CheckPasswordHash("fail", hashedPassword)
	if matchFalse == true {
		t.Error(`[TestCheckPasswordHash] expect matchFalse be false, it was true`)
	}

	matchTrue, _ := auth.CheckPasswordHash(pass, hashedPassword)
	if matchTrue == false {
		t.Errorf(`[TestCheckPasswordHash] expected CheckPasswordHash be true, it was false`)
	}

}
