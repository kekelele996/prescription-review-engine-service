package util

import (
	"testing"
	"time"
)

func TestJWTGenerateAndParse(t *testing.T) {
	secret := "test-secret-abcdef"
	token, err := GenerateToken(secret, time.Hour, 42, "dr_zhang", "doctor")
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if claims.UserID != 42 || claims.Username != "dr_zhang" || claims.Role != "doctor" {
		t.Errorf("claims = %+v", claims)
	}
	if claims.Issuer != Issuer {
		t.Errorf("issuer = %s", claims.Issuer)
	}
}

func TestJWTWrongSecret(t *testing.T) {
	token, err := GenerateToken("secret-a", time.Hour, 1, "u", "admin")
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	if _, err := ParseToken("secret-b", token); err == nil {
		t.Fatal("expected parse error with wrong secret")
	}
}
