package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DemoSecret is a demonstration-only HMAC key. It provides obfuscation, NOT security.
// Do NOT use in production. Replace with a proper asymmetric signing scheme.
const DemoSecret = "asf-enterprise-secret-2024"

// ASFVersion is the current version. Overridden at build time via -ldflags -X.
var ASFVersion = "5.1.2"

const LicensePrefix = "ASF"

type LicenseInfo struct {
	Key     string
	Valid   bool
	Tier    string
	Message string
}

func LoadLicense() *LicenseInfo {
	path := asfLicensePath()
	data, err := os.ReadFile(path)
	if err != nil {
		oldPath := oldLicensePath()
		if oldPath != "" {
			data, err = os.ReadFile(oldPath)
			if err == nil {
				if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
					debugLog.Printf("license migrate mkdir: %v", err)
				} else if err := os.WriteFile(path, data, 0600); err != nil {
					debugLog.Printf("license migrate write: %v", err)
				} else {
					os.Remove(oldPath)
				}
			}
		}
	}
	if err != nil {
		return &LicenseInfo{Valid: false, Message: "No license key found. See https://github.com/moksh5936-2/asfassumption/issues"}
	}

	key := strings.TrimSpace(string(data))
	return ValidateLicense(key)
}

// ValidateLicense checks license key format and signature.
// Supports Ed25519 (ASF-ED25519-*) and legacy HMAC (ASF-XXXX-*) formats.
func ValidateLicense(key string) *LicenseInfo {
	key = strings.TrimSpace(key)

	// Try Ed25519 first
	if strings.HasPrefix(key, "ASF-ED25519-") {
		if ValidateEd25519(key) {
			return &LicenseInfo{
				Key:     key,
				Valid:   true,
				Tier:    "Enterprise",
				Message: fmt.Sprintf("✓ Enterprise License - %s", key),
			}
		}
		return &LicenseInfo{Key: key, Valid: false, Message: "Invalid Ed25519 license signature"}
	}

	// Legacy HMAC format
	parts := strings.Split(key, "-")
	if len(parts) != 5 || parts[0] != LicensePrefix {
		return &LicenseInfo{Key: key, Valid: false, Message: "Invalid license format. Expected: ASF-ED25519-* or ASF-XXXX-XXXX-XXXX-XXXX"}
	}

	rawData := parts[1] + parts[2] + parts[3]
	sig := parts[4]

	mac := hmac.New(sha256.New, []byte(DemoSecret))
	mac.Write([]byte(rawData))
	expected := hex.EncodeToString(mac.Sum(nil))[:8]

	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return &LicenseInfo{Key: key, Valid: false, Message: "Invalid license signature"}
	}

	tier := "Enterprise"
	if strings.Contains(key, "TRIAL") {
		tier = "Trial"
	}

	return &LicenseInfo{
		Key:     key,
		Valid:   true,
		Tier:    tier,
		Message: fmt.Sprintf("✓ %s License - %s", tier, key),
	}
}

func SaveLicense(key string) error {
	path := asfLicensePath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strings.TrimSpace(key)), 0600)
}

func GenerateLicenseKey(data string) (string, error) {
	if strings.HasPrefix(strings.ToUpper(data), "ED25519") {
		return GenerateEd25519License(data)
	}
	mac := hmac.New(sha256.New, []byte(DemoSecret))
	mac.Write([]byte(data))
	sig := hex.EncodeToString(mac.Sum(nil))[:8]
	if len(data) < 12 {
		return "", fmt.Errorf("data must be at least 12 characters for HMAC license")
	}
	parts := []string{data[:4], data[4:8], data[8:12]}
	return fmt.Sprintf("ASF-%s-%s-%s-%s", parts[0], parts[1], parts[2], sig), nil
}

func licenseDataForHMAC(key string) string {
	parts := strings.Split(key, "-")
	if len(parts) < 4 {
		return key
	}
	return parts[1] + "-" + parts[2] + "-" + parts[3]
}
