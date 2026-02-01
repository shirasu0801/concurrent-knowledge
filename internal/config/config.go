package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Port string
	Env  string

	// Database configuration
	DatabasePath string

	// Logging
	LogLevel string

	// Cache configuration
	CacheTTL int

	// Application configuration
	QuestionsPerQuiz  int
	TimePerQuestion   int // seconds
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{
		Port:              getEnv("PORT", "8080"),
		Env:               getEnv("ENV", "development"),
		DatabasePath:      getEnv("DATABASE_PATH", "./data/quiz.db"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		CacheTTL:          getEnvAsInt("CACHE_TTL", 3600),
		QuestionsPerQuiz:  getEnvAsInt("QUESTIONS_PER_QUIZ", 10),
		TimePerQuestion:   getEnvAsInt("TIME_PER_QUESTION", 60),
	}

	log.Printf("Configuration loaded: Port=%s, Env=%s, DB=%s", cfg.Port, cfg.Env, cfg.DatabasePath)

	return cfg
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as int or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}
