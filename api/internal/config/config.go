package config

import (
	"os"
	"strings"
)

// Config holds runtime settings, sourced exclusively from env vars.
type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	AppURL             string
	FrontendURL        string
	AllowedOrigins     []string
	GoogleClientID     string
	GoogleClientSecret string
	YandexClientID     string
	YandexClientSecret string
}

func Load() Config {
	return Config{
		Port:               env("PORT", "8080"),
		DatabaseURL:        env("DATABASE_URL", "postgres://flashcards:flashcards@localhost:5432/flashcards?sslmode=disable"),
		JWTSecret:          env("JWT_SECRET", "dev-secret-change-me"),
		AppURL:             env("APP_URL", "http://localhost:8080"),
		FrontendURL:        env("FRONTEND_URL", "http://localhost:5175"),
		AllowedOrigins:     splitList(env("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:5174,http://localhost:5175,http://localhost:5176")),
		GoogleClientID:     env("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: env("GOOGLE_CLIENT_SECRET", ""),
		YandexClientID:     env("YANDEX_CLIENT_ID", ""),
		YandexClientSecret: env("YANDEX_CLIENT_SECRET", ""),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// splitList parses a comma-separated env value, trimming blanks and trailing slashes.
func splitList(v string) []string {
	var out []string
	for _, part := range strings.Split(v, ",") {
		part = strings.TrimRight(strings.TrimSpace(part), "/")
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
