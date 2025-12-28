package env

import (
	"strconv"
	"time"

	envLoader "github.com/himtar/go-boilerplate/pkg/env"
)

// Variables is an immutable struct representing configuration variables.
type Variables struct {
	env                 string
	dbURI               string
	db                  string
	port                string
	moduleName          string
	readTimeoutMs       *time.Duration
	writeTimeoutMs      *time.Duration
	shutdownTimeoutMs   *time.Duration
	idleTimeoutMs       *time.Duration
	readHeaderTimeoutMs *time.Duration
}

func parseMsDuration(val string) (*time.Duration, error) {
	if val == "" {
		return nil, nil
	}
	ms, err := strconv.Atoi(val)
	if err != nil {
		return nil, err
	}
	d := time.Duration(ms) * time.Millisecond
	return &d, nil
}

// function to load env variables
func LoadENVVariables() (*Variables, error) {
	err := envLoader.LoadENV()
	if err != nil {
		return nil, err
	}

	readTimeoutMs, err := parseMsDuration(envLoader.GetEnvOrDefault("READ_TIMEOUT_MS", ""))
	if err != nil {
		return nil, err
	}
	writeTimeoutMs, err := parseMsDuration(envLoader.GetEnvOrDefault("WRITE_TIMEOUT_MS", ""))
	if err != nil {
		return nil, err
	}
	shutdownTimeoutMs, err := parseMsDuration(envLoader.GetEnvOrDefault("SHUTDOWN_TIMEOUT_MS", ""))
	if err != nil {
		return nil, err
	}
	idleTimeoutMs, err := parseMsDuration(envLoader.GetEnvOrDefault("IDLE_TIMEOUT_MS", ""))
	if err != nil {
		return nil, err
	}
	readHeaderTimeoutMs, err := parseMsDuration(envLoader.GetEnvOrDefault("READ_HEADER_TIMEOUT_MS", ""))
	if err != nil {
		return nil, err
	}

	return &Variables{
		env:                 envLoader.GetEnvOrDefault("ENV", "development"),
		dbURI:               envLoader.GetEnvOrDefault("DB_URI", ""),
		db:                  envLoader.GetEnvOrDefault("DB", ""),
		port:                envLoader.GetEnvOrDefault("PORT", ":8080"),
		moduleName:          envLoader.GetEnvOrDefault("MODULE_NAME", ""),
		readTimeoutMs:       readTimeoutMs,
		writeTimeoutMs:      writeTimeoutMs,
		shutdownTimeoutMs:   shutdownTimeoutMs,
		idleTimeoutMs:       idleTimeoutMs,
		readHeaderTimeoutMs: readHeaderTimeoutMs,
	}, nil
}

// Accessor methods to retrieve the values (no setters provided).

func (v *Variables) Env() string {
	return v.env
}

func (v *Variables) DB_URI() string {
	return v.dbURI
}

func (v *Variables) DB() string {
	return v.db
}

func (v *Variables) Port() string {
	return v.port
}

func (v *Variables) MODULE_NAME() string {
	return v.moduleName
}

func (v *Variables) READ_TIMEOUT_MS() *time.Duration {
	return v.readTimeoutMs
}

func (v *Variables) WRITE_TIMEOUT_MS() *time.Duration {
	return v.writeTimeoutMs
}

func (v *Variables) SHUTDOWN_TIMEOUT_MS() *time.Duration {
	return v.shutdownTimeoutMs
}

func (v *Variables) IDLE_TIMEOUT_MS() *time.Duration {
	return v.idleTimeoutMs
}

func (v *Variables) READ_HEADER_TIMEOUT_MS() *time.Duration {
	return v.readHeaderTimeoutMs
}
