package main

import (
	"os"
	"path/filepath"
	"runtime"
)

func asfConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return filepath.Join(os.TempDir(), "asf", "config")
	}
	switch runtime.GOOS {
	case "linux":
		return filepath.Join(configDir, "asf")
	case "darwin":
		return filepath.Join(configDir, "asf")
	default:
		return filepath.Join(configDir, "ASF")
	}
}

func asfCacheDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return filepath.Join(os.TempDir(), "asf", "cache")
	}
	return filepath.Join(cacheDir, "asf")
}

func asfDataDir() string {
	switch runtime.GOOS {
	case "linux":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".local", "share", "asf")
	case "darwin":
		return asfConfigDir()
	default:
		return asfConfigDir()
	}
}

func asfEngineDir() string {
	return filepath.Join(asfDataDir(), "engine")
}

func asfConfigPath() string {
	return filepath.Join(asfConfigDir(), "config.yaml")
}

func asfLicensePath() string {
	return filepath.Join(asfConfigDir(), "license.key")
}

func asfLogPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(asfCacheDir(), "asf.log")
	}
	return filepath.Join(home, ".asf", "logs", "asf.log")
}

func asfLogsDir() string {
	return filepath.Dir(asfLogPath())
}

func ensureRuntimeDirs() error {
	for _, d := range []string{asfCacheDir(), asfConfigDir(), asfDataDir(), asfEngineDir(), asfLogsDir()} {
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

func oldConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".asf", "config.yaml")
}

func oldLicensePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".asf", "license.key")
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
