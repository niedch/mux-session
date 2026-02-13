package conf

import (
	"path/filepath"

	"github.com/niedch/mux-session/internal/logger"

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
	SearchPaths     []string        `koanf:"search_paths"`
	PreviewProvider *string         `koanf:"preview_provider"`
	Default         ProjectConfig   `koanf:"default"`
	Project         []ProjectConfig `koanf:"project"`
}

func Load(configFile string) (*Config, error) {
	k := koanf.New(".")

	if err := loadProjectConfig(k, configFile); err != nil {
		return nil, err
	}

	var conf Config

	if err := k.UnmarshalWithConf("", &conf, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		logger.Fatalf("Cannot load Config: %v\n", err)
		return nil, err
	}

	if err := validateConfig(&conf); err != nil {
		logger.Fatalf("Config validation failed: %v\n", err)
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
		logger.Fatalf("Cannot load config from '%s': %v\n", configPath, err)
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
