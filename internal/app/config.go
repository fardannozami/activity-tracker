package app

import "os"

type Config struct {
	DatabaseUrl string
	HTTPAddr    string
}

func LoadConfig() (Config, error) {
	cfg := Config{
		DatabaseUrl: getEnv("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/activity_tracker?sslmode=disable"),
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
