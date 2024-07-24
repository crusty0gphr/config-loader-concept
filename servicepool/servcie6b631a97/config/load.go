package config

import (
	_ "embed"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// -------- Base config ----------------

const externalConfigPathEnv = "CONFIG_PATH"

type AppConfig struct {
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
	ConfigPath string `yaml:"config_path"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Base struct {
	App    AppConfig    `yaml:"app"`
	Server ServerConfig `yaml:"server"`
}

//go:embed cfg.base.yaml
var configYaml []byte

func LoadBase() Base {
	var cfg Base

	err := yaml.Unmarshal(configYaml, &cfg)
	if err != nil {
		log.Fatalf("unable to read base config: %v", err)
		return Base{}
	}

	return cfg
}

func (b Base) GetExternalConfigPath() string {
	return getEnv(
		externalConfigPathEnv,
		b.App.ConfigPath,
	)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// -------- External config ----------------

const natsUrlEnv = "NATS_URL"

type External struct {
	Database Database `yaml:"database"`
	NATS     Nats     `yaml:"nats"`
	Logging  Logging  `yaml:"logging"`
	Features Features `yaml:"features"`
}

type Database struct {
	Host string `yaml:"host"`
}

type Nats struct {
	Host string `yaml:"host"`
}

type Logging struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Features struct {
	EnableFeatureX bool `yaml:"enable_feature_x"`
	EnableFeatureY bool `yaml:"enable_feature_y"`
}

func LoadExternal(path string) External {
	var cfg External

	file, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading YAML file: %s\n", err)
		return cfg
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		log.Printf("Error unmarshaling YAML file: %s\n", err)
		return cfg
	}

	return cfg
}

func (cfg External) BuildNatsUrl() string {
	return getEnv(
		natsUrlEnv,
		cfg.NATS.Host,
	)
}
