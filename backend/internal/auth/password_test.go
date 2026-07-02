package auth_test

import (
	"testing"

	"github.com/readr/api/internal/auth"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "supersecret123"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == password {
		t.Error("hash should not equal the plain password")
	}
	if !auth.CheckPassword(hash, password) {
		t.Error("CheckPassword should return true for correct password")
	}
}

func TestWrongPasswordFails(t *testing.T) {
	hash, _ := auth.HashPassword("correctpassword")
	if auth.CheckPassword(hash, "wrongpassword") {
		t.Error("CheckPassword should return false for wrong password")
	}
}
