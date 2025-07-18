package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	Server   ServerConfig
	NLP      NLPConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	Database int
}

// ServerConfig holds server configuration
type ServerConfig struct {
	HTTPPort int
	GRPCPort int
}

// NLPConfig holds NLP service configuration
type NLPConfig struct {
	ServiceURL string
	Timeout    int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{}

	// Database configuration
	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnvAsInt("DB_PORT", 5432)
	config.Database.Name = getEnv("DB_NAME", "ai_context_tracker")
	config.Database.User = getEnv("DB_USER", "tracker_user")
	config.Database.Password = getEnv("DB_PASSWORD", "tracker_password")
	config.Database.SSLMode = getEnv("DB_SSL_MODE", "disable")

	// Redis configuration
	config.Redis.Host = getEnv("REDIS_HOST", "localhost")
	config.Redis.Port = getEnvAsInt("REDIS_PORT", 6379)
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.Database = getEnvAsInt("REDIS_DATABASE", 0)

	// Server configuration
	config.Server.HTTPPort = getEnvAsInt("HTTP_PORT", 8080)
	config.Server.GRPCPort = getEnvAsInt("GRPC_PORT", 9090)

	// NLP service configuration
	config.NLP.ServiceURL = getEnv("NLP_SERVICE_URL", "http://localhost:5000")
	config.NLP.Timeout = getEnvAsInt("NLP_TIMEOUT", 30)

	return config, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// ConnectionString returns the database connection string
func (d *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

// Address returns the Redis address
func (r *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}