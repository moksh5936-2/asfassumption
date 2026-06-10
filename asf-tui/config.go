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
		PythonPath string `yaml:"python_path"`
		ProjectDir string `yaml:"project_dir"`
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
	c.Engine.PythonPath = ""
	c.Engine.ProjectDir = ""
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
			if err := os.MkdirAll(filepath.Dir(newPath), 0755); err == nil {
				data, _ := os.ReadFile(oldPath)
				if data != nil {
					os.WriteFile(newPath, data, 0644)
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
