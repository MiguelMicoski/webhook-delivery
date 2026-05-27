package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8090"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
