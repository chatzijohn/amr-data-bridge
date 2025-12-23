package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type DBConfig struct {
	HOST     string
	PORT     string
	DB       string
	USER     string
	PASSWORD string
}

type ServerConfig struct {
	PORT        string
	HOST        string
	PREFERENCES string
}

type AppConfig struct {
	DB          DBConfig
	ENVIRONMENT string
	TELEMETRY   bool
	SERVER      ServerConfig
}

type ExportPreferences struct {
	WaterMeterFields []string `yaml:"water_meter_fields"`
}

// TokenPolicy defines the rules for a specific API token, like IP whitelisting.
type TokenPolicy struct {
	Enabled bool     `yaml:"enabled"`
	IPs     []string `yaml:"ips"`
}

// AuthConfig holds the authentication configuration loaded from preferences.yaml.
// It defines the header to check and the policies for named tokens.
type AuthConfig struct {
	Header string                 `yaml:"header"`
	Tokens map[string]TokenPolicy `yaml:"tokens"`
}

// AuthTokens is a map of token names to their secret values.
type AuthTokens map[string]string

type Preferences struct {
	Export ExportPreferences `yaml:"export"`
	Auth   AuthConfig        `yaml:"auth"`
}

func Load() *AppConfig {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Application configuration
	environment := os.Getenv("ENVIRONMENT")
	telemetry := strings.ToLower(os.Getenv("TELEMETRY")) == "true"

	// Database configuration
	db := DBConfig{
		HOST:     os.Getenv("POSTGRES_HOST"),
		PORT:     os.Getenv("POSTGRES_PORT"),
		DB:       os.Getenv("POSTGRES_DB"),
		USER:     os.Getenv("POSTGRES_USER"),
		PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
	}

	// Server configuration
	server := ServerConfig{
		HOST:        getEnvWithDefault("HOST", "0.0.0.0"),
		PORT:        getEnvWithDefault("PORT", "8080"),
		PREFERENCES: strings.ToLower(getEnvWithDefault("PREFERENCES_FILE", "./preferences.yaml")),
	}

	return &AppConfig{ENVIRONMENT: environment, TELEMETRY: telemetry, DB: db, SERVER: server}

}

// Helper function to get env with fallback
func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// LoadPreferences loads preferences.yaml from disk.
func LoadPreferences(path string) (*Preferences, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var prefs Preferences
	if err := yaml.Unmarshal(data, &prefs); err != nil {
		return nil, err
	}

	return &prefs, nil
}

// LoadTokens parses a comma-separated string of key:value pairs from the TOKEN_LIST env var.
// e.g., "TOKEN_LIST=name1:value1,name2:value2"
func LoadTokens() (AuthTokens, error) {
	const envVar = "TOKEN_LIST"
	tokens := make(AuthTokens)
	tokenList := os.Getenv(envVar)
	if tokenList == "" {
		return tokens, nil // No tokens defined, return empty map
	}

	pairs := strings.Split(tokenList, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid token format in %s: %s", envVar, pair)
		}
		tokens[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return tokens, nil
}
