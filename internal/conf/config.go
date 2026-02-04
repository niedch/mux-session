package conf

import (
	"log"
	"path/filepath"

	"github.com/adrg/xdg"
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
	Name         *string           `koanf:"name"`
	WindowConfig []WindowConfig    `koanf:"window"`
	Env          map[string]string `koanf:"env"`
}

type Config struct {
	SearchPaths []string        `koanf:"search_paths"`
	Default     ProjectConfig   `koanf:"default"`
	Project     []ProjectConfig `koanf:"project"`
}

func Load(configFile string) (*Config, error) {
	k := koanf.New(".")

	if err := loadProjectConfig(k, configFile); err != nil {
		return nil, err
	}

	var conf Config

	if err := k.UnmarshalWithConf("", &conf, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		log.Fatal("Cannot load Config: ", err)
		return nil, err
	}

	if err := validateConfig(&conf); err != nil {
		log.Fatal("Config validation failed: ", err)
		return nil, err
	}

	return &conf, nil
}

func loadProjectConfig(k *koanf.Koanf, configFile string) error {
	configPath := configFile
	if configPath == "" {
		var err error
		configPath, err = xdg.SearchConfigFile(filepath.Join("mux-session", "config.toml"))
		if err != nil {
			return err
		}
	}

	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		log.Fatal("Cannot load config from '"+configPath+"'", err)
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
	for _, projectConfig := range c.Project {
		if *projectConfig.Name == dir {
			return &projectConfig
		}
	}

	return nil
}
