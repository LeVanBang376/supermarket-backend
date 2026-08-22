package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	AdminUsername  string
	AdminEmail     string
	AdminPassword  string
	AllowedOrigins []string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	cfg := &Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		AdminUsername:  os.Getenv("ADMIN_USERNAME"),
		AdminEmail:     os.Getenv("ADMIN_EMAIL"),
		AdminPassword:  os.Getenv("ADMIN_PASSWORD"),
		AllowedOrigins: strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	return cfg
}
