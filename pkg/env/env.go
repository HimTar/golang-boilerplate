package env

import (
	"os"

	"github.com/joho/godotenv"
)

// function to load env variables
func LoadENV() error {
	// Load variables from the .env file
	err := godotenv.Load()

	return err
}

// getEnvOrDefault retrieves the value of the environment variable or returns the default value.
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
