package config

import "os"

// Config holds runtime settings, sourced exclusively from env vars.
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func Load() Config {
	return Config{
		Port:        env("PORT", "8080"),
		DatabaseURL: env("DATABASE_URL", "postgres://flashcards:flashcards@localhost:5432/flashcards?sslmode=disable"),
		JWTSecret:   env("JWT_SECRET", "dev-secret-change-me"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
