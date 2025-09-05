package app

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const (
	envFile = ".env"
)

type (
	AppConfig struct {
		Database   DatabaseConfig `yaml:"database"`
		Logging    LoggingConfig  `yaml:"logging"`
		ConfigFile string
	}
)

var (
	Database *DatabaseConfig
	Logging  *LoggingConfig
)

func Setup() {
	err := godotenv.Load(envFile)
	if err != nil {
		logrus.Debug(err)
	}

	var Http = &AppConfig{
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
	}
	Http.Database.Setup()
	fmt.Println("*************** APP DATABASE SETUP FINISHED ***************")
	Http.Logging.Setup()
	Database = &Http.Database
	Logging = &Http.Logging
}

func Config(key string) string {
	value := os.Getenv(key)
	return value
}
