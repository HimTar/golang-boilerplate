package libraries

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Variables is an immutable struct representing configuration variables.
type Variables struct {
	env   string
	dbURI string
	db    string
	port  string
}

// function to load env variables
func LoadENVVariables() *Variables {

	// Load variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Variables{
		env:   getEnvOrDefault("ENV", "development"),
		dbURI: getEnvOrDefault("DB_URI", ""),
		db:    getEnvOrDefault("DB", ""),
		port:  getEnvOrDefault("PORT", ":8080"),
	}
}

// getEnvOrDefault retrieves the value of the environment variable or returns the default value.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Accessor methods to retrieve the values (no setters provided).

func (v *Variables) Env() string {
	return v.env
}

func (v *Variables) DBURI() string {
	return v.dbURI
}

func (v *Variables) DB() string {
	return v.db
}

func (v *Variables) Port() string {
	return v.port
}
