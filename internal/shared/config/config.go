package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	AdminAPIKey string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{
		Port:        port,
		DatabaseURL: os.Getenv("DATABASE_URL"),
		AdminAPIKey: os.Getenv("ADMIN_API_KEY"),
	}
}
