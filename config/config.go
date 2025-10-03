package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	HOST     string
	DB       string
	USER     string
	PASSWORD string
}

type ServerConfig struct {
	PORT string
	HOST string
}

type AppConfig struct {
	DB          DBConfig
	ENVIRONMENT string
	SERVER      ServerConfig
}

func Load() *AppConfig {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Application configuration
	environment := os.Getenv("ENVIRONMENT")

	// Database configuration
	db := DBConfig{
		HOST:     os.Getenv("POSTGRES_HOST"),
		DB:       os.Getenv("POSTGRES_DB"),
		USER:     os.Getenv("POSTGRES_USER"),
		PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
	}

	// Server configuration
	server := ServerConfig{
		HOST: getEnvWithDefault("HOST", "127.0.0.1"),
		PORT: getEnvWithDefault("PORT", "8080"),
	}

	return &AppConfig{ENVIRONMENT: environment, DB: db, SERVER: server}

}

// Helper function to get env with fallback
func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
