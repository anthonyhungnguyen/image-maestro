package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	PostgresConfig struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	} `yaml:"postgres"`
}

func readConfig(configPath string) (Config, error) {
	config, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var postgresConfig Config

	err = yaml.Unmarshal(config, &postgresConfig)

	if err != nil {
		return Config{}, err
	}

	return postgresConfig, nil
}

func GetConnection() *sql.DB {
	Config, err := readConfig("static/config.yaml")
	log.Info().Msg(fmt.Sprintf("%+v", Config))

	if err != nil {
		log.Fatal().Err(err).Msg("Error reading config file")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", Config.PostgresConfig.Host, Config.PostgresConfig.Port, Config.PostgresConfig.User, Config.PostgresConfig.Password, Config.PostgresConfig.Dbname)

	// Open database connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening database connection")
	}

	// Validate connection
	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("Error validating database connection")
	}
	return db
}
