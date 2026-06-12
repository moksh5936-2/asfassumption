package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var latestVersionCache struct {
	version string
	checked time.Time
}

const versionCheckURL = "https://api.github.com/repos/moksh5936-2/asfassumption/releases/latest"

type githubRelease struct {
	TagName string `json:"tag_name"`
}

// CheckLatestVersion fetches the latest release version from GitHub.
// Returns empty string if check fails (network, rate limit, etc).
// Results are cached for 1 hour to avoid hammering the API.
func CheckLatestVersion() string {
	if !latestVersionCache.checked.IsZero() && time.Since(latestVersionCache.checked) < time.Hour {
		return latestVersionCache.version
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(versionCheckURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return ""
	}

	v := strings.TrimPrefix(rel.TagName, "v")
	latestVersionCache.version = v
	latestVersionCache.checked = time.Now()
	return v
}

// VersionCheckMessage returns a user-facing message if a newer version exists.
func VersionCheckMessage() string {
	latest := CheckLatestVersion()
	if latest == "" {
		return ""
	}
	current := strings.TrimPrefix(ASFVersion, "v")
	if latest == current {
		return ""
	}
	// Simple comparison: if strings differ, suggest upgrade
	return fmt.Sprintf("A newer version (v%s) is available. You are running v%s.", latest, current)
}
