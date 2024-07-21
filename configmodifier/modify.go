package configmodifier

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"

	"gopkg.in/yaml.v3"
)

type Database struct {
	Host string `yaml:"user"`
}

type Logging struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Features struct {
	EnableFeatureX bool `yaml:"enable_feature_x"`
	EnableFeatureY bool `yaml:"enable_feature_y"`
}

type Config struct {
	Database Database `yaml:"database"`
	Logging  Logging  `yaml:"logging"`
	Features Features `yaml:"features"`
}

type dbCred struct {
	user   string
	dbName string
	pass   string
	port   string
}

type dbCredentials []dbCred

func (d dbCredentials) GetRandom() dbCred {
	return d[rand.IntN(len(d))]
}

var dbCredPool = dbCredentials{
	{user: "fba1977e", dbName: "c3f106287d67", pass: "fba1977e-6359-4292-9405-c3f106287d67", port: "5433"},
	{user: "1f65e1cf", dbName: "5f93ac4a5894", pass: "1f65e1cf-7a57-4e4b-8e1b-5f93ac4a5894", port: "5434"},
	{user: "d84ae75d", dbName: "e18a1046226a", pass: "d84ae75d-f795-4be5-8979-e18a1046226a", port: "5435"},
}

var backup = "configs/cfg.backup.yaml"

func randomModifyConfig(config *Config) {
	dbc := dbCredPool.GetRandom()

	config.Database.Host = fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		dbc.user,
		dbc.pass,
		dbc.port,
		dbc.dbName,
	)
	config.Logging.Level = []string{"DEBUG", "INFO", "WARN", "ERROR"}[rand.IntN(4)]
	config.Features.EnableFeatureX = rand.IntN(2) == 1
	config.Features.EnableFeatureY = rand.IntN(2) == 1
}

func ModifyConfigFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading YAML file: %s\n", err)
		return
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Printf("Error unmarshaling YAML file: %s\n", err)
		return
	}

	randomModifyConfig(&config)

	mod, err := yaml.Marshal(&config)
	if err != nil {
		log.Printf("Error marshaling modified config: %s\n", err)
		return
	}

	err = os.WriteFile(path, mod, os.ModePerm)
	if err != nil {
		log.Printf("Error writing modified YAML file: %s\n", err)
		return
	}

	log.Printf("Successfully modified config file: %s\n", path)
}
