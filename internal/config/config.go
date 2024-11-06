package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment       string `default:"development"`
	Port              string `default:"4000"`
	DatabaseURL       string `default:"empty"`
	SessionCookieName string `default:"session"`
	DatabaseMode      string `default:"none"`
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v.", err)
	}
	return &Config{
		Environment:       os.Getenv("ENV"),
		Port:              os.Getenv("HTTP_PORT"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		SessionCookieName: os.Getenv("SESSION_COOKIE"),
		DatabaseMode:      os.Getenv("DATABASE_MODE"),
	}
}
