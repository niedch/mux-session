package conf

import (
	"log"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type PanelConfig struct {
	Cmd string `koanf:"cmd"`
}

type WindowConfig struct {
	WindowName  string        `koanf:"window_name"`
	PanelConfig []PanelConfig `koanf:"panel"`
	Cmd         *string       `koanf:"cmd"`
}

type ProjectConfig struct {
	Name         *string        `koanf:"name"`
	WindowConfig []WindowConfig `koanf:"window"`
}

type Config struct {
	Search_paths []string        `koanf:"search_paths"`
	Default      ProjectConfig   `koanf:"default"`
	Projects     []ProjectConfig `koanf:"project"`
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

	return &conf, nil
}

func loadProjectConfig(k *koanf.Koanf) error {
	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		log.Fatal("Cannot load config from 'config.toml'", err)
		return err
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
	for _, projectConfig := range c.Projects {
		if *projectConfig.Name == dir {
			return &projectConfig
		}
	}

	return nil
}
