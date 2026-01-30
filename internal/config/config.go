package config

import (
	"errors"
	"os"
)

type Config struct {
	ApiVersion string
	ApiPort    string
	*PostgresConf
}

type PostgresConf struct {
	DbName   string
	Host     string
	Port     string
	User     string
	Password string
}

func NewConfig() (*Config, error) {
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	apiVersion := os.Getenv("API_VERSION")
	apiPort := os.Getenv("API_PORT")

	if len(dbName) == 0 || len(host) == 0 || len(port) == 0 || len(user) == 0 || len(password) == 0 {
		return nil, errors.New("error getting some DB env")
	}

	pConf := &PostgresConf{
		DbName:   dbName,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}

	return &Config{
		ApiVersion:   apiVersion,
		ApiPort:      apiPort,
		PostgresConf: pConf,
	}, nil
}
