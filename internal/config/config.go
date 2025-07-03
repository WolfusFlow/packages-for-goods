package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port string

	JetViewsPath string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string
}

func Load() Config {
	return Config{
		Port: getEnv("PORT", "8080"),

		JetViewsPath: getEnv("JET_VIEWS_PATH", "templates"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "packaging"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret: getEnv("JWT_SECRET", "your-secret"),
	}
}

func (c Config) GetPostgresURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
