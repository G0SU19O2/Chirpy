package auth

import (
	"testing"
	"time"
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

func TestMakeAndValidateJWT(t *testing.T) {
	// Test setup
	userID := "123"
	secret := "test-secret-key"
	expiresIn := time.Hour

	// Test token creation
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}
	if token == "" {
		t.Error("Token should not be empty")
	}

	// Test token validation
	extractedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}
	if extractedUserID != userID {
		t.Errorf("Expected userID %s, got %s", userID, extractedUserID)
	}

	// Test with wrong secret
	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Error("ValidateJWT should fail with wrong secret")
	}

	// Test with invalid token
	_, err = ValidateJWT("invalid-token", secret)
	if err == nil {
		t.Error("ValidateJWT should fail with invalid token")
	}
}
