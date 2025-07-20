package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	Server struct {
		Port string
	}

	Secret struct {
		Secret string
	}
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{}

	cfg.DB.Host = getEnv("DB_HOST", "db")
	cfg.DB.Port = getEnv("DB_PORT", "5432")
	cfg.DB.User = getEnv("DB_USER", "user")
	cfg.DB.Password = getEnv("DB_PASSWORD", "password")
	cfg.DB.Name = getEnv("DB_NAME", "vk")
	cfg.DB.SSLMode = getEnv("DB_SSLMODE", "disable")

	cfg.Server.Port = getEnv("SERVER_PORT", "8080")

	cfg.Secret.Secret = getEnv("SECRET", "")

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
