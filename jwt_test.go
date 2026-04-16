package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"strings"
	"testing"
)

func generateTestKey(t *testing.T) (*rsa.PrivateKey, []byte) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	return key, pemBytes
}

func TestGenerateJWT_Valid(t *testing.T) {
	key, pemBytes := generateTestKey(t)

	jwt, err := GenerateJWT("12345", pemBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		t.Fatalf("JWT should have 3 parts, got %d", len(parts))
	}

	// Verify header
	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("failed to decode header: %v", err)
	}
	var header map[string]string
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		t.Fatalf("failed to parse header: %v", err)
	}
	if header["alg"] != "RS256" {
		t.Errorf("alg = %q, want RS256", header["alg"])
	}
	if header["typ"] != "JWT" {
		t.Errorf("typ = %q, want JWT", header["typ"])
	}

	// Verify claims
	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("failed to decode claims: %v", err)
	}
	var claims map[string]any
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		t.Fatalf("failed to parse claims: %v", err)
	}
	if claims["iss"] != "12345" {
		t.Errorf("iss = %v, want 12345", claims["iss"])
	}
	if _, ok := claims["iat"]; !ok {
		t.Error("missing iat claim")
	}
	if _, ok := claims["exp"]; !ok {
		t.Error("missing exp claim")
	}

	// Verify signature
	sigBytes, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		t.Fatalf("failed to decode signature: %v", err)
	}
	signingInput := parts[0] + "." + parts[1]
	hash := sha256.Sum256([]byte(signingInput))
	if err := rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA256, hash[:], sigBytes); err != nil {
		t.Errorf("signature verification failed: %v", err)
	}
}

func TestGenerateJWT_PKCS8Key(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Bytes,
	})

	jwt, err := GenerateJWT("99999", pemBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if jwt == "" {
		t.Error("JWT should not be empty")
	}
}

func TestGenerateJWT_InvalidPEM(t *testing.T) {
	_, err := GenerateJWT("12345", []byte("not a pem"))
	if err == nil {
		t.Error("expected error for invalid PEM")
	}
}
