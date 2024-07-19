package config

import (
	_ "embed"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

type App struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Nats struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type ConfigsList map[string]string

type Config struct {
	App     App         `yaml:"app"`
	Nats    Nats        `yaml:"nats"`
	Configs ConfigsList `yaml:"configs"`
}

//go:embed cfg.base.yaml
var configYaml []byte

func Load() Config {
	var cfg Config

	err := yaml.Unmarshal(configYaml, &cfg)
	if err != nil {
		log.Fatalf("unable to read base config: %v", err)
		return Config{}
	}

	return cfg
}

func (cfg Config) BuildNatsUrl() string {
	return fmt.Sprintf("nats://%s:%s", cfg.Nats.Host, cfg.Nats.Port)
}
