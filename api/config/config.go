package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	HTTPPort  string
	HTTPSPort string

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret                     string
	AccessTokenExpireDuration  time.Duration
	RefreshTokenExpireDuration time.Duration
}

func Load() *Config {
	return &Config{
		Server:   *LoadServerConfig(),
		Database: *LoadDBConfig(),
		JWT:      *LoadJWTConfig(),
	}
}

func LoadServerConfig() *ServerConfig {
	return &ServerConfig{
		HTTPPort:     getEnv("SERVER_HTTP_PORT", "8080"),
		HTTPSPort:    getEnv("SERVER_HTTPS_PORT", "8443"),
		ReadTimeout:  getDuration("SERVER_READ_TIMEOUT", 30*time.Second),
		WriteTimeout: getDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
	}
}

func LoadDBConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("POSTGRES_HOST", ""),
		Port:     getEnv("POSTGRES_PORT", ""),
		Username: getEnv("POSTGRES_USER", ""),
		Password: getEnv("POSTGRES_PASSWORD", ""),
		Name:     getEnv("POSTGRES_DB", ""),
		SSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),
	}
}

func LoadJWTConfig() *JWTConfig {
	return &JWTConfig{
		Secret:                     getEnv("JWT_SECRET", ""),
		AccessTokenExpireDuration:  getDuration("JWT_ACCESS_EXPIRE_DURATION", 15*time.Minute),
		RefreshTokenExpireDuration: getDuration("JWT_REFRESH_EXPIRE_DURATION", 7*24*time.Hour),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
