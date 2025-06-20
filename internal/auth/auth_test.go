package auth

import (
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "supersecret123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == password {
		t.Error("Hash should not be the same as the password")
	}

	if err := CheckPassword(password, hash); err != nil {
		t.Errorf("CheckPassword failed for correct password: %v", err)
	}

	wrongPassword := "wrongPassword"
	if err := CheckPassword(wrongPassword, hash); err == nil {
		t.Error("CheckPassword should fail for incorrect password")
	}
}
