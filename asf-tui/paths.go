package main

import (
	"os"
	"path/filepath"
	"runtime"
)

func asfRootDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), "asf")
	}
	return filepath.Join(home, ".asf")
}

func asfConfigDir() string {
	return filepath.Join(asfRootDir(), "config")
}

func asfCacheDir() string {
	return filepath.Join(asfRootDir(), "cache")
}

func asfDataDir() string {
	return filepath.Join(asfRootDir(), "data")
}

func asfEngineDir() string {
	return filepath.Join(asfDataDir(), "engine")
}

func asfConfigPath() string {
	return filepath.Join(asfConfigDir(), "config.json")
}

func asfLicensePath() string {
	return filepath.Join(asfRootDir(), "license.key")
}

func asfLogPath() string {
	return filepath.Join(asfRootDir(), "logs", "asf.log")
}

func asfLogsDir() string {
	return filepath.Dir(asfLogPath())
}

func asfTelemetryPath() string {
	return filepath.Join(asfConfigDir(), "telemetry.json")
}

func ensureRuntimeDirs() error {
	for _, d := range []string{asfCacheDir(), asfConfigDir(), asfDataDir(), asfEngineDir(), asfLogsDir(), asfRootDir()} {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	return nil
}

func legacyConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "asf", "config.yaml")
}

func xdgConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	switch runtime.GOOS {
	case "linux", "darwin":
		return filepath.Join(configDir, "asf", "config.yaml")
	default:
		return filepath.Join(configDir, "ASF", "config.yaml")
	}
}

func oldLicensePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".asf", "license.key")
}

func oldLicenseAtRootPath() string {
	return filepath.Join(asfRootDir(), "license.key")
}

func ensureDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}

func getDefaultPath(cfg *Config) (string, error) {
	if cfg != nil && cfg.Output.Directory != "" {
		return cfg.Output.Directory, nil
	}
	return os.Getwd()
}
