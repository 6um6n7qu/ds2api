// Package config provides configuration management for ds2api.
// It loads settings from environment variables with sensible defaults.
package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	// Server settings
	ServerHost string
	ServerPort string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// DS2 / Synology NAS settings
	NASHost     string
	NASPort     string
	NASUser     string
	NASPassword string
	NASProtocol string

	// API settings
	APIKey      string
	DebugMode   bool
	LogLevel    string
}

// Load reads configuration from environment variables.
// Missing required variables will be left empty; callers should validate.
func Load() *Config {
	return &Config{
		ServerHost:   getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		ReadTimeout:  getDurationEnv("READ_TIMEOUT_SEC", 60) * time.Second,
		WriteTimeout: getDurationEnv("WRITE_TIMEOUT_SEC", 60) * time.Second,

		NASHost:     getEnv("NAS_HOST", ""),
		NASPort:     getEnv("NAS_PORT", "5000"),
		NASUser:     getEnv("NAS_USER", ""),
		NASPassword: getEnv("NAS_PASSWORD", ""),
		NASProtocol: getEnv("NAS_PROTOCOL", "http"),

		APIKey:    getEnv("API_KEY", ""),
		DebugMode: getBoolEnv("DEBUG", false),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
	}
}

// Validate checks that all required configuration fields are set.
// Returns a slice of error messages; an empty slice means the config is valid.
func (c *Config) Validate() []string {
	var errs []string

	if c.NASHost == "" {
		errs = append(errs, "NAS_HOST is required")
	}
	if c.NASUser == "" {
		errs = append(errs, "NAS_USER is required")
	}
	if c.NASPassword == "" {
		errs = append(errs, "NAS_PASSWORD is required")
	}

	return errs
}

// Address returns the full listen address in host:port format.
func (c *Config) Address() string {
	return c.ServerHost + ":" + c.ServerPort
}

// NASBaseURL returns the base URL for the Synology NAS API.
func (c *Config) NASBaseURL() string {
	return c.NASProtocol + "://" + c.NASHost + ":" + c.NASPort
}

// getEnv returns the value of the environment variable named by key,
// or fallback if the variable is not set or empty.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getBoolEnv parses a boolean environment variable.
func getBoolEnv(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

// getDurationEnv parses an integer environment variable as a time.Duration multiplier.
func getDurationEnv(key string, fallback int) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return time.Duration(fallback)
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return time.Duration(fallback)
	}
	return time.Duration(n)
}
