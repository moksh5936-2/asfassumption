package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	General struct {
		Theme    string `yaml:"theme"`
		FoxStyle string `yaml:"fox_style"`
	} `yaml:"general"`
	Analysis struct {
		Depth         string `yaml:"depth"`
		Stride        bool   `yaml:"stride"`
		Controls      bool   `yaml:"controls"`
		RiskThreshold string `yaml:"risk_threshold"`
	} `yaml:"analysis"`
	AI struct {
		Enabled         bool     `yaml:"enabled"`
		ActiveModel     string   `yaml:"active_model"`
		InstalledModels []string `yaml:"installed_models"`
	} `yaml:"ai"`
	Output struct {
		Default   string `yaml:"default"`
		Directory string `yaml:"directory"`
	} `yaml:"output"`
	Appearance struct {
		Theme    string `yaml:"theme"`
		FoxStyle string `yaml:"fox_style"`
	} `yaml:"appearance"`
	Engine struct {
		UseNativeEngine bool `yaml:"use_native_engine"`
	} `yaml:"engine"`
}

func DefaultConfig() Config {
	c := Config{}
	c.General.Theme = "Dark"
	c.General.FoxStyle = "Classic"
	c.Analysis.Depth = "deep"
	c.Analysis.Stride = true
	c.Analysis.Controls = true
	c.Analysis.RiskThreshold = "low"
	c.AI.Enabled = false
	c.AI.ActiveModel = ""
	c.AI.InstalledModels = []string{}
	c.Output.Default = "markdown"
	c.Output.Directory = "./reports"
	c.Appearance.Theme = "Dark"
	c.Appearance.FoxStyle = "Classic"
	c.Engine.UseNativeEngine = true
	return c
}

func ConfigPath() string {
	newPath := asfConfigPath()
	oldPath := legacyConfigPath()
	if oldPath == "" {
		return newPath
	}
	if _, err := os.Stat(oldPath); err == nil {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
				debugLog.Printf("migrate mkdir: %v", err)
			} else {
				data, err := os.ReadFile(oldPath)
				if err != nil {
					debugLog.Printf("migrate read: %v", err)
				} else {
					if err := os.WriteFile(newPath, data, 0644); err != nil {
						debugLog.Printf("migrate write: %v", err)
					}
				}
			}
		}
	}
	return newPath
}

func LoadConfig(path string) (*Config, error) {
	c := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return &c, nil
			}
			if writeErr := c.Save(path); writeErr != nil {
				debugLog.Printf("failed to write default config: %v", writeErr)
			}
			return &c, nil
		}
		return &c, err
	}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return &c, err
	}
	return &c, nil
}

func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
