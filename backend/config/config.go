package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration values loaded from environment variables.
type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

// LoadConfig reads the .env file and returns a populated Config struct.
// The application will terminate if the .env file cannot be loaded.
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file, using system environment variables")
	}

	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "appointment_booking"),
		JWTSecret:  getEnv("JWT_SECRET", "default-secret-key"),
	}
}

// getEnv retrieves the value of an environment variable or returns the fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	} 
	return fallback
}
