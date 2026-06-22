package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	JWTExpiryHours int
	ServerPort     string
}

func Load() *Config {
	return &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBUser:         getEnv("DB_USER", "priceuser"),
		DBPassword:     getEnv("DB_PASSWORD", "pricepass"),
		DBName:         getEnv("DB_NAME", "pricecompare"),
		JWTSecret:      getEnv("JWT_SECRET", "secret"),
		JWTExpiryHours: getEnvInt("JWT_EXPIRY_HOURS", 72),
		ServerPort:     getEnv("BACKEND_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
