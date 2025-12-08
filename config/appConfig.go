package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort string
	Dsn        string // data source name
	AppSecret  string
}

func SetupEnv() (cfg AppConfig, err error) {
	// Load .env (ignore error because file may not exist in production)
	_ = godotenv.Load(".env")

	cfg.ServerPort = os.Getenv("HTTP_PORT")
	if cfg.ServerPort == "" {
		return cfg, errors.New("HTTP_PORT not found")
	}

	cfg.Dsn = os.Getenv("DSN")
	if cfg.Dsn == "" {
		return cfg, errors.New("DSN not found")
	}

	cfg.AppSecret = os.Getenv("APP_SECRET")
	if cfg.AppSecret == "" {
		return cfg, errors.New("APP_SECRET not found")
	}

	return cfg, nil
}
