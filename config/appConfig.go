package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort            string
	Dsn                   string
	AppSecret             string
	TwilioAccountSid      string
	TwilioAuthToken       string
	TwilioFromPhoneNumber string
	StripeSecret          string
	Pubkey                string
	SuccessUrl            string
	CancelUrl             string
}

func SetupEnv() (cfg AppConfig, err error) {

	// Load .env (ignore errors, because in production env variables come from the OS)
	_ = godotenv.Load(".env")

	cfg = AppConfig{
		ServerPort:            os.Getenv("HTTP_PORT"),
		Dsn:                   os.Getenv("DSN"),
		AppSecret:             os.Getenv("APP_SECRET"),
		TwilioAccountSid:      os.Getenv("TWILIO_ACCOUNT_SID"),
		TwilioAuthToken:       os.Getenv("TWILIO_AUTH_TOKEN"),
		TwilioFromPhoneNumber: os.Getenv("TWILIO_FROM_PHONE_NUMBER"),
		StripeSecret:          os.Getenv("STRIPE_SECRET"),
		Pubkey:                os.Getenv("STRIPE_PUB_KEY"),
		SuccessUrl:            os.Getenv("SUCCESS_URL"),
		CancelUrl:             os.Getenv("CANCEL_URL"),
	}

	// Validate required variables
	if cfg.ServerPort == "" {
		return cfg, errors.New("environment variable HTTP_PORT is missing")
	}
	if cfg.Dsn == "" {
		return cfg, errors.New("environment variable DSN is missing")
	}
	if cfg.AppSecret == "" {
		return cfg, errors.New("environment variable APP_SECRET is missing")
	}
	if cfg.TwilioAccountSid == "" {
		return cfg, errors.New("environment variable TWILIO_ACCOUNT_SID is missing")
	}
	if cfg.TwilioAuthToken == "" {
		return cfg, errors.New("environment variable TWILIO_AUTH_TOKEN is missing")
	}
	if cfg.TwilioFromPhoneNumber == "" {
		return cfg, errors.New("environment variable TWILIO_FROM_PHONE_NUMBER is missing")
	}

	return cfg, nil
}
