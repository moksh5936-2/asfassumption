package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

var (
	demoPrivateKey ed25519.PrivateKey
	demoPublicKey  ed25519.PublicKey
)

func init() {
	seed := sha256.Sum256([]byte("asf-ed25519-demo-seed-2024-ed25519"))
	demoPrivateKey = ed25519.NewKeyFromSeed(seed[:])
	demoPublicKey = demoPrivateKey.Public().(ed25519.PublicKey)
}

func Ed25519PublicKeyHex() string {
	return hex.EncodeToString(demoPublicKey)
}

func ValidateEd25519(key string) bool {
	key = strings.TrimSpace(key)
	if !strings.HasPrefix(key, "ASF-ED25519-") {
		return false
	}
	parts := strings.SplitN(key, "-", 4)
	if len(parts) != 4 {
		return false
	}
	sigHex := parts[3]
	sig, err := hex.DecodeString(sigHex)
	if err != nil || len(sig) != ed25519.SignatureSize {
		return false
	}
	msg := []byte(strings.Join(parts[:3], "-"))
	return ed25519.Verify(demoPublicKey, msg, sig)
}

func GenerateEd25519License(identifier string) (string, error) {
	if len(identifier) < 4 {
		return "", fmt.Errorf("identifier must be at least 4 characters")
	}
	clean := strings.ToUpper(strings.ReplaceAll(identifier, "-", ""))
	payload := fmt.Sprintf("ASF-ED25519-%s", clean)
	msg := []byte(payload)
	sig := ed25519.Sign(demoPrivateKey, msg)
	return fmt.Sprintf("%s-%s", payload, hex.EncodeToString(sig)), nil
}

func ReplacePublicKey(pub ed25519.PublicKey) bool {
	if len(pub) != ed25519.PublicKeySize {
		return false
	}
	demoPublicKey = pub
	return true
}
