package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	JWTSecret          string
	SuperAdminUsername string
	SuperAdminEmail    string
	SuperAdminPassword string
	SuperAdminFullName string
	SuperAdminPhone    string
	AllowedOrigins     []string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	cfg := &Config{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		SuperAdminUsername: os.Getenv("SUPER_ADMIN_USERNAME"),
		SuperAdminEmail:    os.Getenv("SUPER_ADMIN_EMAIL"),
		SuperAdminPassword: os.Getenv("SUPER_ADMIN_PASSWORD"),
		SuperAdminFullName: os.Getenv("SUPER_ADMIN_FULL_NAME"),
		SuperAdminPhone:    os.Getenv("SUPER_ADMIN_PHONE"),
		AllowedOrigins:     strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return cfg
}
