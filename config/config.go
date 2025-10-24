package config

import (
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

type Preferences struct {
	Export ExportPreferences `yaml:"export"`
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
