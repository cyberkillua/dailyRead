package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	// Add other configuration parameters as needed
}

func LoadConfig() (*Config, error) {
	portString := os.Getenv("PORT")
	if portString == "" {
		return nil, fmt.Errorf("PORT environment variable not set")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	return &Config{
		Port:        portString,
		DatabaseURL: dbUrl,
	}, nil
}
