package auth_test

import (
	"testing"

	"github.com/readr/api/internal/auth"
)

var (
	accessSecret  = []byte("test-access-secret")
	refreshSecret = []byte("test-refresh-secret")
)

func TestGenerateAndValidateTokenPair(t *testing.T) {
	userID := "507f1f77bcf86cd799439011"

	pair, err := auth.GenerateTokenPair(userID, accessSecret, refreshSecret)
	if err != nil {
		t.Fatalf("GenerateTokenPair: %v", err)
	}

	claims, err := auth.ValidateToken(pair.AccessToken, accessSecret)
	if err != nil {
		t.Fatalf("ValidateToken (access): %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, claims.UserID)
	}
	if claims.Type != auth.TokenTypeAccess {
		t.Errorf("expected type %s, got %s", auth.TokenTypeAccess, claims.Type)
	}

	claims, err = auth.ValidateToken(pair.RefreshToken, refreshSecret)
	if err != nil {
		t.Fatalf("ValidateToken (refresh): %v", err)
	}
	if claims.Type != auth.TokenTypeRefresh {
		t.Errorf("expected type %s, got %s", auth.TokenTypeRefresh, claims.Type)
	}
}

func TestAccessTokenCannotBeUsedAsRefresh(t *testing.T) {
	pair, _ := auth.GenerateTokenPair("user1", accessSecret, refreshSecret)

	// Access token signed with accessSecret should fail validation against refreshSecret.
	_, err := auth.ValidateToken(pair.AccessToken, refreshSecret)
	if err == nil {
		t.Error("expected error when validating access token with refresh secret")
	}
}

func TestInvalidTokenIsRejected(t *testing.T) {
	_, err := auth.ValidateToken("not.a.real.token", accessSecret)
	if err == nil {
		t.Error("expected error for malformed token")
	}
}
