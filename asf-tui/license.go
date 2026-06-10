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

const (
	ASFVersion    = "1.0.0"
	LicensePrefix = "ASF"
)

type LicenseInfo struct {
	Key       string
	Valid     bool
	Tier      string
	Message   string
}

func LoadLicense() *LicenseInfo {
	home, err := os.UserHomeDir()
	if err != nil {
		return &LicenseInfo{Valid: false, Message: "Cannot determine home directory"}
	}

	path := filepath.Join(home, ".asf", "license.key")
	data, err := os.ReadFile(path)
	if err != nil {
		return &LicenseInfo{Valid: false, Message: "No license key found. Contact support@asfsecurity.com"}
	}

	key := strings.TrimSpace(string(data))
	return ValidateLicense(key)
}

func ValidateLicense(key string) *LicenseInfo {
	key = strings.TrimSpace(key)
	parts := strings.Split(key, "-")
	if len(parts) != 5 || parts[0] != LicensePrefix {
		return &LicenseInfo{Key: key, Valid: false, Message: "Invalid license format. Expected: ASF-XXXX-XXXX-XXXX-XXXX"}
	}

	rawData := parts[1] + parts[2] + parts[3]
	sig := parts[4]

	mac := hmac.New(sha256.New, []byte("asf-enterprise-secret-2024"))
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
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, ".asf")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "license.key"), []byte(strings.TrimSpace(key)), 0600)
}

func GenerateLicenseKey(data string) string {
	mac := hmac.New(sha256.New, []byte("asf-enterprise-secret-2024"))
	mac.Write([]byte(data))
	sig := hex.EncodeToString(mac.Sum(nil))[:8]
	parts := []string{data[:4], data[4:8], data[8:12]}
	return fmt.Sprintf("ASF-%s-%s-%s-%s", parts[0], parts[1], parts[2], sig)
}

func licenseDataForHMAC(key string) string {
	parts := strings.Split(key, "-")
	if len(parts) < 4 {
		return key
	}
	return parts[1] + "-" + parts[2] + "-" + parts[3]
}
