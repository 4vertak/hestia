package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() (*Config, error) {
	return load(".env")
}

func load(path string) (*Config, error) {
	values, err := godotenv.Read(path)
	if err != nil && !os.IsNotExist(err) {
		// Parser errors can contain secret values from the file.
		return nil, fmt.Errorf("cannot load configuration file %q", path)
	}
	get := func(key, fallback string) string {
		value, exists := os.LookupEnv(key)
		if !exists {
			value = values[key]
		}
		if value == "" {
			return fallback
		}
		return value
	}
	cfg := &Config{
		ServerPort: get("SERVER_PORT", "8080"),
		DBHost:     get("DB_HOST", "127.0.0.1"),
		DBPort:     get("DB_PORT", "5432"),
		DBUser:     get("DB_USER", "hestia"),
		DBPassword: get("DB_PASSWORD", ""),
		DBName:     get("DB_NAME", "store"),
		DBSSLMode:  get("DB_SSLMODE", "disable"),
	}
	if cfg.DBPassword == "" {
		return nil, errors.New("DB_PASSWORD must be set in the environment or .env")
	}
	return cfg, nil
}
