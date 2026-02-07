package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret     string
	JWTExpiration string

	// Gemini AI
	GeminiAPIKey string
}

var AppConfig *Config

func Load() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		// Server
		ServerPort: getEnv("SERVER_PORT", "8080"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "smart_food_store"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// JWT
		JWTSecret:     getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),

		// Gemini AI
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
	}

	return AppConfig
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetDSN() string {
	return "host=" + c.DBHost +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" port=" + c.DBPort +
		" sslmode=" + c.DBSSLMode
}
