package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds all ASF0 user settings, persisted as JSON.
type Config struct {
	General struct {
		Theme    string `json:"theme" yaml:"theme"`
		FoxStyle string `json:"fox_style" yaml:"fox_style"`
		Debug    bool   `json:"debug" yaml:"debug"`
	} `json:"general" yaml:"general"`
	Analysis struct {
		Depth         string `json:"depth" yaml:"depth"`
		Stride        bool   `json:"stride" yaml:"stride"`
		Controls      bool   `json:"controls" yaml:"controls"`
		RiskThreshold string `json:"risk_threshold" yaml:"risk_threshold"`
	} `json:"analysis" yaml:"analysis"`
	AI struct {
		Enabled         bool     `json:"enabled" yaml:"enabled"`
		ActiveModel     string   `json:"active_model" yaml:"active_model"`
		InstalledModels []string `json:"installed_models" yaml:"installed_models"`
	} `json:"ai" yaml:"ai"`
	Output struct {
		Default   string `json:"default" yaml:"default"`
		Directory string `json:"directory" yaml:"directory"`
	} `json:"output" yaml:"output"`
	Appearance struct {
		Theme    string `json:"theme" yaml:"theme"`
		FoxStyle string `json:"fox_style" yaml:"fox_style"`
	} `json:"appearance" yaml:"appearance"`
	Engine struct {
		UseNativeEngine bool `json:"use_native_engine" yaml:"use_native_engine"`
	} `json:"engine" yaml:"engine"`
	Telemetry struct {
		OptIn bool   `json:"opt_in" yaml:"opt_in"`
		ID    string `json:"id,omitempty" yaml:"id,omitempty"`
	} `json:"telemetry,omitempty" yaml:"telemetry,omitempty"`
	Encryption struct {
		Enabled  bool   `json:"enabled" yaml:"enabled"`
		KeyPath  string `json:"key_path,omitempty" yaml:"key_path,omitempty"`
		KeyCheck string `json:"-" yaml:"-"`
	} `json:"encryption,omitempty" yaml:"encryption,omitempty"`
	StrictCLI bool `json:"strict_cli,omitempty" yaml:"strict_cli,omitempty"`
	Version   int  `json:"_version,omitempty" yaml:"_version,omitempty"`
}

const configVersion = 2

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	c := Config{
		Version: configVersion,
	}
	c.General.Theme = "ASF0"
	c.General.FoxStyle = "Classic"
	c.General.Debug = false
	c.Analysis.Depth = "deep"
	c.Analysis.Stride = true
	c.Analysis.Controls = true
	c.Analysis.RiskThreshold = "low"
	c.AI.Enabled = false
	c.AI.ActiveModel = ""
	c.AI.InstalledModels = []string{}
	c.Output.Default = "markdown"
	c.Output.Directory = "./reports"
	c.Appearance.Theme = "ASF0"
	c.Appearance.FoxStyle = "Classic"
	c.Engine.UseNativeEngine = true
	c.Telemetry.OptIn = false
	c.Encryption.Enabled = false
	return c
}

func ConfigPath() string {
	return asfConfigPath()
}

func migrateLegacyConfigs() error {
	newPath := ConfigPath()
	if _, err := os.Stat(newPath); err == nil {
		return nil
	}

	var legacyPaths []string
	if p := legacyConfigPath(); p != "" {
		legacyPaths = append(legacyPaths, p)
	}
	if p := xdgConfigPath(); p != "" && p != legacyConfigPath() {
		legacyPaths = append(legacyPaths, p)
	}

	for _, oldPath := range legacyPaths {
		if _, err := os.Stat(oldPath); os.IsNotExist(err) {
			continue
		}
		debugLog.Printf("migrating config from %s", oldPath)
		data, err := os.ReadFile(oldPath)
		if err != nil {
			debugLog.Printf("migrate read %s: %v", oldPath, err)
			continue
		}
		var legacy Config
		if yaml.Unmarshal(data, &legacy) == nil {
			legacy.Version = configVersion
			if err := SaveConfig(newPath, &legacy); err != nil {
				debugLog.Printf("migrate write %s: %v", newPath, err)
			} else {
				debugLog.Printf("migrated config from %s to %s", oldPath, newPath)
				return nil
			}
		}
	}
	return nil
}

func LoadConfig(path string) (*Config, error) {
	if err := migrateLegacyConfigs(); err != nil {
		debugLog.Printf("config migration: %v", err)
	}

	c := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if mkErr := os.MkdirAll(filepath.Dir(path), 0755); mkErr != nil {
				return &c, nil
			}
			if writeErr := SaveConfig(path, &c); writeErr != nil {
				debugLog.Printf("failed to write default config: %v", writeErr)
			}
			return &c, nil
		}
		return &c, err
	}

	if json.Unmarshal(data, &c) == nil {
		if c.Version < configVersion {
			c.Version = configVersion
			SaveConfig(path, &c)
		}
		return &c, nil
	}

	if yaml.Unmarshal(data, &c) == nil {
		c.Version = configVersion
		if saveErr := SaveConfig(path, &c); saveErr != nil {
			debugLog.Printf("migrate yaml->json save: %v", saveErr)
		}
		return &c, nil
	}

	return &c, fmt.Errorf("config file %s: unable to parse as JSON or YAML", path)
}

func SaveConfig(path string, c *Config) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (c *Config) Save(path string) error {
	return SaveConfig(path, c)
}
