package conf

import (
	"errors"
	"log"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type PanelConfig struct {
	PanelDirection string `koanf:"panel_direction"`
	Cmd            string `koanf:"cmd"`
}

type WindowConfig struct {
	WindowName  string        `koanf:"window_name"`
	PanelConfig []PanelConfig `koanf:"panel_config"`
	Primary     *bool         `koanf:"primary"`
	Cmd         *string       `koanf:"cmd"`
}

type ProjectConfig struct {
	Name         *string        `koanf:"name"`
	WindowConfig []WindowConfig `koanf:"window"`
}

type Config struct {
	SearchPaths []string        `koanf:"search_paths"`
	Default     ProjectConfig   `koanf:"default"`
	Project     []ProjectConfig `koanf:"project"`
}

func Load() (*Config, error) {
	k := koanf.New(".")

	if err := loadProjectConfig(k); err != nil {
		return nil, err
	}

	var conf Config

	if err := k.UnmarshalWithConf("", &conf, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		log.Fatal("Cannot load Config", err)
		return nil, err
	}

	if err := validateConfig(&conf); err != nil {
		log.Fatal("Config validation failed", err)
		return nil, err
	}

	return &conf, nil
}

func loadProjectConfig(k *koanf.Koanf) error {
	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		log.Fatal("Cannot load config from 'config.toml'", err)
		return err
	}
	return nil
}

func validateConfig(conf *Config) error {
	for _, project := range conf.Project {
		if err := validateProjectConfig(project); err != nil {
			return err
		}
	}

	return validateProjectConfig(conf.Default)
}

func validateProjectConfig(project ProjectConfig) error {
	if err := validatePrimaryMarker(project); err != nil {
		return err
	}

	if err := validatePanelConfig(project); err != nil {
		return err
	}

	return nil
}

func validatePrimaryMarker(project ProjectConfig) error {
	primaryCount := 0
	for _, window := range project.WindowConfig {
		// Check if this window is marked as primary
		if window.Primary != nil && *window.Primary {
			primaryCount++
		}

		// Validate that only one window is marked as primary
		if primaryCount > 1 {
			return errors.New("only one window can be marked as primary in project configuration")
		}
	}

	return nil
}

func validatePanelConfig(project ProjectConfig) error {
	for _, window := range project.WindowConfig {
		for _, panel := range window.PanelConfig {
			if panel.PanelDirection != "v" && panel.PanelDirection != "h" {
				return errors.New("panel_direction must be 'v' or 'h'")
			}
		}
	}

	return nil
}

// Finds Project Config otherwise returns Default
func (c *Config) GetProjectConfig(dir string) ProjectConfig {
	projectConfig := c.findProject(dir)

	if projectConfig == nil {
		return c.Default
	}

	return *projectConfig
}

func (c *Config) findProject(dir string) *ProjectConfig {
	for _, projectConfig := range c.Project {
		if *projectConfig.Name == dir {
			return &projectConfig
		}
	}

	return nil
}
