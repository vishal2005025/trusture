package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
		Host string
		Mode string // debug, release, test
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	JWT struct {
		Secret     string
		ExpiryHours int
	}
	Blockchain struct {
		PolygonRPC    string
		PrivateKey    string
		GasLimit      int64
		GasPriceGwei  int64
	}
	Platform struct {
		FeePercentage float64
		Environment   string
	}
	Logging struct {
		Level  string
		Format string // json, text
	}
}

var AppConfig *Config

func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{}

	// Server configuration
	config.Server.Port = getEnv("SERVER_PORT", "8080")
	config.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
	config.Server.Mode = getEnv("GIN_MODE", "debug")

	// Database configuration
	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnv("DB_PORT", "5432")
	config.Database.User = getEnv("DB_USER", "trusture")
	config.Database.Password = getEnv("DB_PASSWORD", "trusture123")
	config.Database.DBName = getEnv("DB_NAME", "trusture_db")
	config.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// JWT configuration
	config.JWT.Secret = getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production")
	config.JWT.ExpiryHours = getEnvInt("JWT_EXPIRY_HOURS", 24)

	// Blockchain configuration
	config.Blockchain.PolygonRPC = getEnv("POLYGON_RPC", "https://polygon-mumbai.g.alchemy.com/v2/demo")
	config.Blockchain.PrivateKey = getEnv("POLYGON_PRIVATE_KEY", "1111111111111111111111111111111111111111111111111111111111111111")
	config.Blockchain.GasLimit = getEnvInt64("POLYGON_GAS_LIMIT", 300000)
	config.Blockchain.GasPriceGwei = getEnvInt64("POLYGON_GAS_PRICE_GWEI", 30)

	// Platform configuration
	config.Platform.FeePercentage = getEnvFloat("PLATFORM_FEE_PERCENTAGE", 1.0)
	config.Platform.Environment = getEnv("ENVIRONMENT", "development")

	// Logging configuration
	config.Logging.Level = getEnv("LOG_LEVEL", "info")
	config.Logging.Format = getEnv("LOG_FORMAT", "json")

	AppConfig = config
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func IsDevelopment() bool {
	return AppConfig.Platform.Environment == "development"
}

func IsProduction() bool {
	return AppConfig.Platform.Environment == "production"
}