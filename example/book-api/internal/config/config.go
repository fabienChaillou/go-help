package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration.
type Config struct {
	DB   DBConfig
	HTTP HTTPConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type HTTPConfig struct {
	Port string
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

// Load reads configuration from environment variables with sane defaults.
func Load() Config {
	return Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "bookuser"),
			Password: getEnv("DB_PASSWORD", "bookpass"),
			Name:     getEnv("DB_NAME", "bookdb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		HTTP: HTTPConfig{
			Port: getEnv("HTTP_PORT", "8080"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
