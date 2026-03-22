package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DBURL        string
	JWTSecret    string
	JWTExpiresIn string

	// Resend email service
	ResendAPIKey string
	ResendFrom   string
	AppBaseURL   string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	// PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// DATABASE
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is missing in .env")
	}

	// JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is missing in .env")
	}

	jwtExp := os.Getenv("JWT_EXPIRES_IN")
	if jwtExp == "" {
		jwtExp = "24h"
	}

	// EMAIL (Resend)
	resendKey := os.Getenv("RESEND_API_KEY")
	if resendKey == "" {
		log.Fatal("RESEND_API_KEY is missing in .env")
	}

	resendFrom := os.Getenv("RESEND_FROM_EMAIL")
	if resendFrom == "" {
		log.Fatal("RESEND_FROM_EMAIL is missing in .env")
	}

	appURL := os.Getenv("APP_BASE_URL")
	if appURL == "" {
		log.Fatal("APP_BASE_URL is missing in .env")
	}

	return &Config{
		Port:         port,
		DBURL:        dbURL,
		JWTSecret:    jwtSecret,
		JWTExpiresIn: jwtExp,
		ResendAPIKey: resendKey,
		ResendFrom:   resendFrom,
		AppBaseURL:   appURL,
	}
}
