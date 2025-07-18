package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port string

	Production bool

	JetViewsPath string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string
	JWTExpiry time.Duration

	AdminEmail    string
	AdminPassword string
}

func Load() *Config {
	expiry, _ := time.ParseDuration(os.Getenv("JWT_EXPIRY")) // e.g. "30m"

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPass := os.Getenv("ADMIN_PASSWORD")

	if adminEmail == "" || adminPass == "" {
		panic("missing required admin credentials: ADMIN_EMAIL and ADMIN_PASSWORD must be set")
	}

	production, err := strconv.ParseBool(getEnv("PRODUCTION", "false"))
	if err != nil {
		panic(fmt.Sprintf("problem parsing production env variable: %v", err))
	}

	return &Config{
		Port: getEnv("PORT", "8080"),

		Production: production,

		JetViewsPath: getEnv("JET_VIEWS_PATH", "internal/templates"),

		DBHost:     getEnv("DB_HOST", "postgres"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "packaging"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret: getEnv("JWT_SECRET", "super-secret-key"),
		JWTExpiry: expiry,

		AdminEmail:    adminEmail,
		AdminPassword: adminPass,
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
