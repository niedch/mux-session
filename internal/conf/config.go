package conf

import (
	"log"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type WindowConfig struct {
	SessionName string `koanf:"session_name"`
}

type ProjectConfig struct {
	Name         *string      `koanf:"name"`
	Dir          *string       `koanf:"dir"`
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
